package providers

import (
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/models"
)

func ytsGetPopular(config *config.Config, logger *zerolog.Logger, r *resty.Client) []models.Film {
	var result YtsPopularRootObject
	_, err := r.R().SetResult(&result).SetHeader("Content-Type", "application/json").Get(config.YtsProvider.PopularUrl)

	if err != nil {
		logger.Fatal().Msgf("Error while retrieving YTS Popular films: %v", err)
	}

	return translate2Model(result)
}

func translate2Model(ytsObject YtsPopularRootObject) []models.Film {
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
