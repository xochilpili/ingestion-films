package providers

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog"
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/models"
)

type Letterbox struct {
	config *config.Config
	logger *zerolog.Logger
	c      *colly.Collector
}

func NewLetterBox(config *config.Config, logger *zerolog.Logger) *Letterbox {
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"),
	)

	return &Letterbox{
		config: config,
		logger: logger,
		c:      c,
	}
}

func (provider *Letterbox) GetFestivals() []models.Film {
	provider.c.Limit(&colly.LimitRule{DomainGlob: "", Parallelism: 2, RandomDelay: time.Duration(provider.config.DelaySecs) * time.Second})

	var items []models.Film

	provider.c.OnHTML(provider.config.LetterboxProvider.FestivalsSelectorRe, func(h *colly.HTMLElement) {
		link := h.Attr("href")
		re := regexp.MustCompile(`festiville/list`)

		if re.MatchString(link) {
			if provider.config.Debug {
				title := h.Text
				provider.logger.Info().Msgf("Link found: %s, Festival: %s", link, title)
			}
			h.Request.Visit(link)
		}
	})

	provider.c.OnHTML("li.poster-container", func(h *colly.HTMLElement) {
		id := h.ChildAttr("div.poster", "data-film-id")
		slug := regexp.MustCompile(`-\d{4}(-\d+)?$`).ReplaceAllString(h.ChildAttr("div.poster", "data-film-slug"), "")
		title := strings.Join(strings.Split(slug, "-"), " ")
		idUrlPath := strings.Join(strings.Split(id, ""), "/")
		imageUrl := fmt.Sprintf("https://a.ltrbxd.com/resized/film-poster/%s/%s-%s-0-125-0-187-crop.jpg", idUrlPath, id, slug)
		provider.logger.Info().Msgf("slug: %s, image: %s", slug, imageUrl)
		items = append(items, models.Film{
			Id:       id,
			Title:    title,
			ImageUrl: imageUrl,
			Year:     time.Now().Year(),
		})
	})

	if provider.config.Debug {
		provider.c.OnResponse(func(r *colly.Response) {
			provider.logger.Info().Msgf("Response: %s", string(r.Body))
		})
	}

	if provider.config.Debug {
		provider.c.OnRequest(func(r *colly.Request) {
			provider.logger.Info().Msgf("Request to: %s", r.URL.String())
		})
	}

	provider.c.Visit(provider.config.LetterboxProvider.FestivalsUrl)
	provider.c.Wait()
	return items
}

func (provider *Letterbox) GetPopular() []models.Film {
	provider.c.Limit(&colly.LimitRule{DomainGlob: "", Parallelism: 2, RandomDelay: 10 * time.Second})

	var films []models.Film
	provider.c.OnHTML(provider.config.LetterboxProvider.PopularSelectorRe, func(h *colly.HTMLElement) {
		id := h.ChildAttr("div.really-lazy-load", "data-film-id")
		slug := h.ChildAttr("div.really-lazy-load", "data-film-slug")
		idUrlPath := strings.Join(strings.Split(id, ""), "/")
		title := h.ChildAttr("img", "alt")
		imageUrl := fmt.Sprintf("https://a.ltrbxd.com/resized/film-poster/%s/%s-%s-0-140-0-210-crop.jpg", idUrlPath, id, slug)
		films = append(films, models.Film{
			Id:       id,
			Title:    title,
			Year:     time.Now().Year(),
			ImageUrl: imageUrl,
		})
	})

	if provider.config.Debug {
		provider.c.OnResponse(func(r *colly.Response) {
			provider.logger.Info().Msgf("Response: %s", string(r.Body))
		})
	}

	if provider.config.Debug {
		provider.c.OnRequest(func(r *colly.Request) {
			provider.logger.Info().Msgf("Requesting URL: %s", r.URL.String())
		})
	}

	provider.c.Visit(fmt.Sprintf("%s%d%s", provider.config.LetterboxProvider.PopularUrl, time.Now().Year(), "/?esiAllowFilters=true"))
	provider.c.Wait()
	return films
}
