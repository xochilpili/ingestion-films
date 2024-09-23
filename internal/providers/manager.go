package providers

import (
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog"
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/models"
)

type GetFestivals func(config *config.Config, logger *zerolog.Logger, c *colly.Collector, r *resty.Client) []models.Film
type GetPopular func(config *config.Config, logger *zerolog.Logger, c *colly.Collector, r *resty.Client) []models.Film
type Handler = struct {
	Enabled      bool
	GetFestivals GetFestivals
	GetPopular   GetPopular
}

type Manager struct {
	config   *config.Config
	logger   *zerolog.Logger
	c        *colly.Collector
	r        *resty.Client
	handlers map[string]Handler
}

func New(config *config.Config, logger *zerolog.Logger) *Manager {
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
		colly.CacheDir("./cache"),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"),
	)

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

	return &Manager{
		config:   config,
		logger:   logger,
		c:        c,
		r:        r,
		handlers: handlers,
	}
}

func (m *Manager) GetFestivals(provider string) []models.Film {
	wg := &sync.WaitGroup{}

	var films []models.Film
	filmsChan := make(chan []models.Film)

	wg.Add(1)

	go func(wg *sync.WaitGroup, filmsChan chan<- []models.Film) {
		defer wg.Done()
		items := m.handlers[provider].GetFestivals(m.config, m.logger, m.c, m.r)
		filmsChan <- items
	}(wg, filmsChan)

	go func() {
		wg.Wait()
		close(filmsChan)
	}()

	for item := range filmsChan {
		films = append(films, item...)
	}

	return films
}

func (m *Manager) GetPopular(provider string) []models.Film {
	wg := &sync.WaitGroup{}

	var films []models.Film
	filmsChan := make(chan []models.Film)

	wg.Add(1)

	for p := range m.handlers {
		if p == provider {
			go func(wg *sync.WaitGroup, filmsChan chan<- []models.Film) {
				defer wg.Done()
				items := m.handlers[provider].GetPopular(m.config, m.logger, m.c, m.r)
				filmsChan <- items
			}(wg, filmsChan)
		}
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
