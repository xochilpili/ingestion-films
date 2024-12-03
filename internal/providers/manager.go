package providers

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/database"
	"github.com/xochilpili/ingestion-films/internal/models"
	"github.com/xochilpili/ingestion-films/internal/services"
	"github.com/xochilpili/ingestion-films/internal/utils"
)

type PgService interface {
	Connect() error
	InsertFilm(table string, columns []string, item *models.Film) error
	GetGenres(genreIds []int) ([]string, error)
	Ping() error
	Close() error
}

type ApiService interface {
	GetMovieDetails(title string) ([]models.TmdbItem, error)
	PlexItemExists(title string) (bool, error)
}

type ProviderConfig struct {
	Enabled              bool
	BaseUrl              string
	PopularUrl           string
	PopularSelectorRe    string
	Festivals            map[string]string
	FestivalsUrl         string
	FestivalsSelectoreRe string
	Debug                bool
	UserAgent            string
	DelaySecs            int
	RequireTmdb          bool
	TmdbUrl              string
	TmdbApiKey           string
	ExcludeGenres        []string
}

type GetFestivals func(config *ProviderConfig, logger *zerolog.Logger, r *resty.Client) []models.Film
type GetPopular func(config *ProviderConfig, logger *zerolog.Logger, r *resty.Client) []models.Film
type PostProcess func(config *ProviderConfig, dbService PgService, apiService ApiService, items *[]models.Film) (*[]models.Film, error)
type Handler = struct {
	Config       *ProviderConfig
	GetFestivals GetFestivals
	GetPopular   GetPopular
	PostProcess  PostProcess
}

type Manager struct {
	config     *config.Config
	logger     *zerolog.Logger
	r          *resty.Client
	handlers   map[string]Handler
	pgService  PgService
	apiService ApiService
}

func New(config *config.Config, logger *zerolog.Logger) *Manager {

	r := resty.New()

	handlers := map[string]Handler{
		"yts": {
			Config: &ProviderConfig{
				Enabled:              true,
				BaseUrl:              "",
				PopularUrl:           "https://yts.mx/api/v2/list_movies.json",
				PopularSelectorRe:    "",
				Festivals:            map[string]string{},
				FestivalsUrl:         "",
				FestivalsSelectoreRe: "",
				Debug:                false,
				UserAgent:            "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
				DelaySecs:            10,
				RequireTmdb:          false,
				TmdbUrl:              "",
				TmdbApiKey:           "",
				ExcludeGenres:        config.ExcludeGenres,
			},
			GetFestivals: nil,
			GetPopular:   ytsGetPopular,
			PostProcess:  nil,
		},
		"letterbox": {
			Config: &ProviderConfig{
				Enabled:              true,
				BaseUrl:              "https://letterboxd.com",
				PopularUrl:           "https://letterboxd.com/films/ajax/popular/this/week/year/",
				PopularSelectorRe:    "li.listitem",
				Festivals:            map[string]string{},
				FestivalsUrl:         "https://letterboxd.com/festiville/lists/",
				FestivalsSelectoreRe: "h2.title-2 a[href]",
				Debug:                false,
				UserAgent:            "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
				DelaySecs:            10,
				RequireTmdb:          true,
				TmdbUrl:              config.Tmdb.Url,
				TmdbApiKey:           config.Tmdb.ApiKey,
				ExcludeGenres:        config.ExcludeGenres,
			},
			GetFestivals: letterboxGetFestivals,
			GetPopular:   letterboxGetPopular,
			PostProcess:  letterboxPostProcess,
		},
		"imdb": {
			Config: &ProviderConfig{
				Enabled:           true,
				BaseUrl:           "",
				PopularUrl:        "https://www.imdb.com/chart/moviemeter/?ref_=nv_mv_mpm",
				PopularSelectorRe: "itemListElement",
				Festivals: map[string]string{
					"cannes":    "https://www.imdb.com/event/ev0000147/",
					"tiff":      "https://www.imdb.com/event/ev0000659/",
					"venecia":   "https://www.imdb.com/event/ev0000681/",
					"oscar":     fmt.Sprintf("https://www.imdb.com/event/ev0000003/%d/", time.Now().Year()),
					"berlinale": "https://www.imdb.com/event/ev0000091/",
				},
				FestivalsUrl:         "",
				FestivalsSelectoreRe: `IMDbReactWidgets\.NomineesWidget\.push\(\[.*?,({.*?})\]\)`,
				Debug:                false,
				UserAgent:            "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
				DelaySecs:            10,
				RequireTmdb:          true,
				TmdbUrl:              config.Tmdb.Url,
				TmdbApiKey:           config.Tmdb.ApiKey,
				ExcludeGenres:        config.ExcludeGenres,
			},
			GetFestivals: imdbGetFestivals,
			GetPopular:   imdbGetPopular,
			PostProcess:  imdbPostProcess,
		},
	}

	pgService := database.New(config, logger)
	pgService.Connect()
	apiService := services.NewApi(config, logger)
	return &Manager{
		config:     config,
		logger:     logger,
		r:          r,
		handlers:   handlers,
		pgService:  pgService,
		apiService: apiService,
	}
}

func (m *Manager) GetFestivals(provider string) []models.Film {
	// yts has no festivals
	if provider == "yts" {
		return nil
	}

	films := m.getFestivals(provider)
	return films

}

func (m *Manager) GetPopular(provider string) []models.Film {
	films := m.getPopular(provider)
	return films
}

func (m *Manager) SyncFestivals(provider string) error {
	if provider == "yts" {
		return errors.New("yts provider has no festivals")
	}
	err := m.pgService.Connect()
	if err != nil {
		m.logger.Err(err).Msg("error while connecting to db")
		return err
	}
	err = m.pgService.Ping()
	if err != nil {
		m.logger.Err(err).Msg("error db do not pong")
		return err
	}

	defer m.pgService.Close()

	films := m.getFestivals(provider)
	for _, item := range films {
		if utils.ExcludeGenre(item.Genre, m.config.ExcludeGenres) {
			continue
		}
		err := m.pgService.InsertFilm("films_festivals", []string{"provider", "title", "year", "genres"}, &item)
		if err != nil {
			m.logger.Err(err).Msgf("error while inserting film %s", item.Title)
			return err
		}
	}
	m.logger.Info().Msgf("sync completed with %d items", len(films))
	return nil
}

func (m *Manager) SyncPopular(provider string) error {
	err := m.pgService.Connect()
	if err != nil {
		m.logger.Fatal().Err(err).Msg("error while connecting to db")
		return err
	}
	err = m.pgService.Ping()
	if err != nil {
		m.logger.Fatal().Err(err).Msg("database didn't pong")
		return err
	}
	defer m.pgService.Close()

	films := m.GetPopular(provider)
	filtered := 0
	for _, item := range films {
		ok, _ := m.apiService.PlexItemExists(item.Title)
		if ok {
			m.logger.Warn().Msgf("film %s already exists in plex", item.Title)
			continue
		}

		if utils.ExcludeGenre(item.Genre, m.config.ExcludeGenres) {
			continue
		}

		err := m.pgService.InsertFilm("films_popular", []string{"provider", "title", "year", "genres"}, &item)
		if err != nil {
			m.logger.Err(err).Msgf("error while inserting film: %s", item.Title)
			return err
		}
		filtered += 1
	}
	m.logger.Info().Msgf("sync completed with %d filtered items from %d", filtered, len(films))
	return nil
}

func (m *Manager) getFestivals(provider string) []models.Film {
	wg := &sync.WaitGroup{}

	var films []models.Film
	filmsChan := make(chan []models.Film)

	for p := range m.handlers {
		if provider != "" && provider != p || p == "yts" {
			continue
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, filmsChan chan<- []models.Film, provider string) {
			defer wg.Done()
			items := m.handlers[provider].GetFestivals(m.handlers[provider].Config, m.logger, m.r)
			m.logger.Info().Msgf("received %d festival films from %s provider", len(items), provider)
			if m.handlers[provider].Config.RequireTmdb {
				m.logger.Info().Msgf("festivals post process for provider: %s", provider)
				m.handlers[provider].PostProcess(m.handlers[provider].Config, m.pgService, m.apiService, &items)
			}
			filmsChan <- items
		}(wg, filmsChan, p)
	}

	go func() {
		wg.Wait()
		close(filmsChan)
	}()

	for item := range filmsChan {
		films = append(films, item...)
	}

	return films
}

func (m *Manager) getPopular(provider string) []models.Film {
	wg := &sync.WaitGroup{}

	var films []models.Film
	filmsChan := make(chan []models.Film)

	for p := range m.handlers {
		if provider != "" && p != provider {
			continue
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, filmsChan chan<- []models.Film, provider string) {
			defer wg.Done()
			items := m.handlers[provider].GetPopular(m.handlers[provider].Config, m.logger, m.r)
			m.logger.Info().Msgf("received %d popular films from %s provider", len(items), provider)
			if m.handlers[provider].Config.RequireTmdb {
				m.logger.Info().Msgf("popular post process for provider: %s", provider)
				m.handlers[provider].PostProcess(m.handlers[provider].Config, m.pgService, m.apiService, &items)
			}
			filmsChan <- items
		}(wg, filmsChan, p)
	}

	go func() {
		wg.Wait()
		close(filmsChan)
	}()

	for items := range filmsChan {
		films = append(films, items...)
	}

	return films
}
