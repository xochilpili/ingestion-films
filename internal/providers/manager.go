package providers

import (
	"errors"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/database"
	"github.com/xochilpili/ingestion-films/internal/models"
)

type PgService interface {
	Connect() error
	InsertFilm(table string, columns []string, item *models.FilmItem) error
	Ping() error
	Close() error
}

type GetFestivals func(config *config.Config, logger *zerolog.Logger, r *resty.Client) []models.Film
type GetPopular func(config *config.Config, logger *zerolog.Logger, r *resty.Client) []models.Film
type Handler = struct {
	Enabled      bool
	GetFestivals GetFestivals
	GetPopular   GetPopular
}

type Manager struct {
	config *config.Config
	logger *zerolog.Logger
	// c         *colly.Collector
	r         *resty.Client
	handlers  map[string]Handler
	pgService PgService
}

func New(config *config.Config, logger *zerolog.Logger) *Manager {

	r := resty.New()

	handlers := map[string]Handler{
		"yts": {
			Enabled:      config.YtsProvider.Enabled,
			GetFestivals: nil,
			GetPopular:   ytsGetPopular,
		},
		"letterbox": {
			Enabled:      config.LetterboxProvider.Enabled,
			GetFestivals: letterboxGetFestivals,
			GetPopular:   letterboxGetPopular,
		},
		"imdb": {
			Enabled:      config.ImdbProvider.Enabled,
			GetFestivals: imdbGetFestivals,
			GetPopular:   imdbGetPopular,
		},
	}

	pgService := database.New(config, logger)

	return &Manager{
		config: config,
		logger: logger,
		// c:         c,
		r:         r,
		handlers:  handlers,
		pgService: pgService,
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
		err := m.pgService.InsertFilm("films_festivals", []string{"provider", "title", "year"}, &models.FilmItem{Provider: item.Provider, Title: item.Title, Year: item.Year})
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
	for _, item := range films {
		err := m.pgService.InsertFilm("films_popular", []string{"provider", "title", "year"}, &models.FilmItem{Provider: item.Provider, Title: item.Title, Year: item.Year})
		if err != nil {
			m.logger.Err(err).Msgf("error while inserting film: %s", item.Title)
			return err
		}
	}
	m.logger.Info().Msgf("sync completed with %d items", len(films))
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
			m.logger.Info().Msgf("festivals from provider: %s", provider)
			items := m.handlers[provider].GetFestivals(m.config, m.logger, m.r)
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
			items := m.handlers[provider].GetPopular(m.config, m.logger, m.r)
			m.logger.Info().Msgf("received %d populae films from %s provider", len(items), provider)
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
