package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog"
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/models"
)

func imdbGetFestivals(config *config.Config, logger *zerolog.Logger, c *colly.Collector, _ *resty.Client) []models.Film {
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1, RandomDelay: time.Duration(config.DelaySecs) * time.Second})

	festivals := config.ImdbProvider.Festivals
	imdbFestivalSelectorRe := `IMDbReactWidgets\.NomineesWidget\.push\(\[.*?,({.*?})\]\)`

	/*
		festivals := map[string]string{
			"cannes":    "https://www.imdb.com/event/ev0000147/",
			"tiff":      "https://www.imdb.com/event/ev0000659/",
			"venecia":   "https://www.imdb.com/event/ev0000681/",
			"oscar":     "https://www.imdb.com/event/ev0000003/",
			"berlinale": "https://www.imdb.com/event/ev0000091/",
		}
	*/

	var model ImdbFestivalRootObject
	var films []models.Film

	logger.Info().Msgf("Festivals: %v\n", festivals)

	for k, url := range festivals {

		c.OnHTML("script", func(h *colly.HTMLElement) {
			scriptContent := h.Text
			re := regexp.MustCompile(imdbFestivalSelectorRe)
			matches := re.FindStringSubmatch(scriptContent)

			if config.Debug {
				fmt.Printf("ScriptContent: %s, matchesLen: %d, matches: %v", scriptContent, len(matches), matches)
				logger.Info().Msgf("ScriptContent: %s, matchesLen: %d, matches: %v", scriptContent, len(matches), matches)
			}

			if len(matches) > 1 {
				// Extract and print the JavaScript object
				jsObject := matches[1]

				err := json.Unmarshal([]byte(jsObject), &model)
				if err != nil {
					log.Fatalf("error while unmarshal %v\n", err)
				}
				films = append(films, translate2FestivalModel(&model)...)
			}
		})

		if config.Debug {
			c.OnRequest(func(r *colly.Request) {
				fmt.Printf("Requesting festival: %s, %s\n", k, r.URL.String())
			})
		}

		logger.Info().Msgf("Visiting Festival: %s, url: %s", k, config.ImdbProvider.HttpPrefix+url)

		c.Visit(config.ImdbProvider.HttpPrefix + url)
		c.Wait()
	}

	return films
}

func imdbGetPopular(config *config.Config, logger *zerolog.Logger, c *colly.Collector, _ *resty.Client) []models.Film {
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1, RandomDelay: time.Duration(config.DelaySecs) * time.Second})

	imdbPopularUrl := config.ImdbProvider.PopularUrl
	imdbPopularSelectorRe := config.ImdbProvider.PopularSelectorRe

	var model ImdbPopularRootObject

	c.OnHTML("script", func(h *colly.HTMLElement) {
		scriptContent := h.Text

		re := regexp.MustCompile(imdbPopularSelectorRe)
		matches := re.FindStringSubmatch(scriptContent)

		if config.Debug {
			logger.Info().Msgf("Captured content: %s, matchesLength: %d, matches: %v", scriptContent, len(matches), matches)
		}

		if len(matches) >= 1 {
			// Extract and print the JavaScript object
			jsObject := scriptContent

			err := json.Unmarshal([]byte(jsObject), &model)
			if err != nil {
				log.Fatalf("error while unmarshal %v\n", err)
			}
		}
	})

	if config.Debug {
		c.OnRequest(func(r *colly.Request) {
			logger.Info().Msgf("Visiting Popular: %s", r.URL.String())
		})
	}

	if config.Debug {
		logger.Info().Msgf("Visiting: %s", imdbPopularUrl)
	}

	c.Visit(imdbPopularUrl)
	c.Wait()
	return translate2PopularModel(&model)
}

func translate2FestivalModel(imdbObject *ImdbFestivalRootObject) []models.Film {
	var films []models.Film
	for _, item := range imdbObject.NomineesWidgetModel.EventEditionSummary.Awards {
		for _, category := range item.Categories {
			for _, nominations := range category.Nominations {
				for _, firstNominee := range nominations.PrimaryNominees {
					films = append(films, models.Film{
						Id:       firstNominee.Const,
						Title:    firstNominee.Name,
						Year:     imdbObject.NomineesWidgetModel.EventEditionSummary.Year,
						ImageUrl: firstNominee.ImageUrl,
					})
				}
			}
		}
	}
	return films
}

func translate2PopularModel(imdbObject *ImdbPopularRootObject) []models.Film {
	var films []models.Film
	for _, film := range imdbObject.ItemListElement {
		var Id string
		parsedUrl, err := url.Parse(film.Item.Url)
		if err != nil {
			Id = ""
		}
		Id = path.Base(parsedUrl.Path)
		genres := strings.Split(film.Item.Genre, ", ")
		films = append(films, models.Film{
			Id:          Id,
			Title:       film.Item.Name,
			Description: film.Item.Description,
			ImageUrl:    film.Item.Image,
			Genre:       genres,
		})
	}
	return films
}
