package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/models"
)

type ApiService struct {
	config *config.Config
	logger *zerolog.Logger
	r      *resty.Client
}

func NewApi(config *config.Config, logger *zerolog.Logger) *ApiService {
	r := resty.New()
	return &ApiService{
		config: config,
		logger: logger,
		r:      r,
	}
}

func (a *ApiService) GetMovieDetails(title string) ([]models.TmdbItem, error) {
	var items models.TmdbResponse
	res, err := a.r.R().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + a.config.Tmdb.ApiKey,
			"User-Agent":    "",
		}).
		SetQueryParams(map[string]string{
			"include_adult": "false",
			"language":      "en-US",
			"page":          "1",
			"query":         title,
		}).
		SetDebug(a.config.Debug).
		Get(a.config.Tmdb.Url + "3/search/movie")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res.Body(), &items)
	if err != nil {
		a.logger.Err(err).Msgf("error while unmarshal response from tmdb: %v", err)
		return nil, err
	}

	return items.Results, nil
}

func (a *ApiService) PlexItemExists(title string) (bool, error) {
	var response models.PlexMediaSearch
	url := fmt.Sprintf("%slibrary/search", a.config.Plex.ApiUrl)
	res, err := a.r.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).
		SetQueryParams(map[string]string{
			"query":     title,
			"sectionId": "1",
			"limit":     "10",
		}).
		SetDebug(a.config.Debug).
		Get(url)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return false, err
	}

	for _, item := range response.MediaContainer.SearchResult {
		if item.Metadata.Title != "" {
			if strings.EqualFold(item.Metadata.Title, title) {
				return true, nil
			}
		}
	}
	return false, nil
}
