package providers

import (
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/models"
)

type YTS struct {
	config  *config.Config
	logger  *zerolog.Logger
	rClient *resty.Client
}

func NewYts(config *config.Config, logger *zerolog.Logger) *YTS {
	r := resty.New()
	return &YTS{
		config:  config,
		logger:  logger,
		rClient: r,
	}
}

func (provider *YTS) GetPopular() []models.Film {
	var result YtsPopularRootObject
	_, err := provider.rClient.R().SetResult(&result).SetHeader("Content-Type", "application/json").SetQueryParams(provider.config.YtsProvider.PopularFilters).Get(provider.config.YtsProvider.PopularUrl)

	if err != nil {
		provider.logger.Fatal().Msgf("Error while retrieving YTS Popular films: %v", err)
	}

	return provider.translate2Model(result)
}

func (provider *YTS) translate2Model(ytsObject YtsPopularRootObject) []models.Film {
	var films []models.Film
	for _, item := range ytsObject.Data.Movies {
		films = append(films, models.Film{
			Id:          item.ImdbCode,
			Title:       item.TitleEnglish,
			Description: item.Summary,
			ImageUrl:    item.BackgroundImageOriginal,
			Genre:       item.Genres,
		})
	}
	return films
}
