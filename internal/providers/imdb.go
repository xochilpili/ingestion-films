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

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog"
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/models"
)

type ImdbProvider struct {
	config *config.Config
	logger *zerolog.Logger
	c      *colly.Collector
}

func NewImdb(config *config.Config, logger *zerolog.Logger) *ImdbProvider {
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(),
	)

	return &ImdbProvider{
		config: config,
		logger: logger,
		c:      c,
	}
}

func (provider *ImdbProvider) GetFestivals() []models.FestivalFilm {
	provider.c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1, RandomDelay: time.Duration(provider.config.ImdbProvider.DelaySecs) * time.Second})

	festivals := provider.config.ImdbProvider.Festivals
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
	var films []models.FestivalFilm

	provider.logger.Info().Msgf("Festivals: %v\n", festivals)

	for k, url := range festivals {

		provider.c.OnHTML("script", func(h *colly.HTMLElement) {
			scriptContent := h.Text
			re := regexp.MustCompile(imdbFestivalSelectorRe)
			matches := re.FindStringSubmatch(scriptContent)

			if provider.config.Debug {
				fmt.Printf("ScriptContent: %s, matchesLen: %d, matches: %v", scriptContent, len(matches), matches)
				provider.logger.Info().Msgf("ScriptContent: %s, matchesLen: %d, matches: %v", scriptContent, len(matches), matches)
			}

			if len(matches) > 1 {
				// Extract and print the JavaScript object
				jsObject := matches[1]

				err := json.Unmarshal([]byte(jsObject), &model)
				if err != nil {
					log.Fatalf("error while unmarshal %v\n", err)
				}
				films = append(films, provider.translate2FestivalModel(&model)...)
			}
		})

		if provider.config.Debug {
			provider.c.OnRequest(func(r *colly.Request) {
				fmt.Printf("Requesting festival: %s, %s\n", k, r.URL.String())
			})
		}

		provider.logger.Info().Msgf("Visiting Festival: %s, url: %s", k, provider.config.ImdbProvider.HttpPrefix+url)

		provider.c.Visit(provider.config.ImdbProvider.HttpPrefix + url)
		provider.c.Wait()
	}

	return films
}

func (provider *ImdbProvider) GetPopular() []models.PopularFilm {
	provider.c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1, RandomDelay: time.Duration(provider.config.ImdbProvider.DelaySecs) * time.Second})

	imdbPopularUrl := provider.config.ImdbProvider.PopularUrl
	imdbPopularSelectorRe := provider.config.ImdbProvider.PopularSelectorRe

	var model ImdbPopularRootObject

	provider.c.OnHTML("script", func(h *colly.HTMLElement) {
		scriptContent := h.Text

		re := regexp.MustCompile(imdbPopularSelectorRe)
		matches := re.FindStringSubmatch(scriptContent)

		if provider.config.Debug {
			provider.logger.Info().Msgf("Captured content: %s, matchesLength: %d, matches: %v", scriptContent, len(matches), matches)
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

	if provider.config.Debug {
		provider.c.OnRequest(func(r *colly.Request) {
			provider.logger.Info().Msgf("Visiting Popular: %s", r.URL.String())
		})
	}

	if provider.config.Debug {
		provider.logger.Info().Msgf("Visiting: %s", imdbPopularUrl)
	}

	provider.c.Visit(imdbPopularUrl)
	provider.c.Wait()
	return provider.translate2PopularModel(&model)
}

func (provider *ImdbProvider) translate2FestivalModel(imdbObject *ImdbFestivalRootObject) []models.FestivalFilm {
	var films []models.FestivalFilm
	for _, item := range imdbObject.NomineesWidgetModel.EventEditionSummary.Awards {
		for _, category := range item.Categories {
			for _, nominations := range category.Nominations {
				for _, firstNominee := range nominations.PrimaryNominees {
					films = append(films, models.FestivalFilm{
						Id:            firstNominee.Const,
						Title:         firstNominee.Name,
						OriginalTitle: firstNominee.OriginalName,
						Year:          imdbObject.NomineesWidgetModel.EventEditionSummary.Year,
						ImageUrl:      firstNominee.ImageUrl,
						Festival:      item.Id,
					})
				}
			}
		}
	}
	return films
}

func (provider *ImdbProvider) translate2PopularModel(imdbObject *ImdbPopularRootObject) []models.PopularFilm {
	var films []models.PopularFilm
	for _, film := range imdbObject.ItemListElement {
		var Id string
		parsedUrl, err := url.Parse(film.Item.Url)
		if err != nil {
			Id = ""
		}
		Id = path.Base(parsedUrl.Path)
		genres := strings.Split(film.Item.Genre, ", ")
		films = append(films, models.PopularFilm{Id: Id, Title: film.Item.Name, Description: film.Item.Description, ImageUrl: film.Item.Image, Genre: genres})
	}
	return films
}
