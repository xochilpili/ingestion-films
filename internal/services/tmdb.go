package services

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/models"
)

type PgService interface {
	Connect() error
	InsertFilm(table string, columns []string, item *models.Film) error
	GetGenres(genreIds []int) ([]string, error)
	Ping() error
	Close() error
}

type TheMovieDatabase struct {
	config    *config.Config
	logger    *zerolog.Logger
	r         *resty.Client
	dbService PgService
}

func NewTmDb(config *config.Config, logger *zerolog.Logger, dbService PgService) *TheMovieDatabase {
	r := resty.New()
	return &TheMovieDatabase{
		config:    config,
		logger:    logger,
		r:         r,
		dbService: dbService,
	}
}

func (s *TheMovieDatabase) GetMovieDetails(ctx context.Context, title string) ([]models.TmdbItem, error) {
	var items models.TmdbResponse
	res, err := s.r.R().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + s.config.Tmdb.ApiKey,
			"User-Agent":    "",
		}).
		SetQueryParams(map[string]string{
			"include_adult": "false",
			"language":      "en-US",
			"page":          "1",
			"query":         title,
		}).
		SetDebug(s.config.Debug).
		Get(s.config.Tmdb.Url + "3/search/movie")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res.Body(), &items)
	if err != nil {
		s.logger.Err(err).Msgf("error while unmarshal response from tmdb: %v", err)
		return nil, err
	}

	return items.Results, nil
}

func (s *TheMovieDatabase) GenresLookup(ids []int) ([]string, error) {
	genres, err := s.dbService.GetGenres(ids)
	if err != nil {
		return nil, err
	}
	return genres, nil
}
