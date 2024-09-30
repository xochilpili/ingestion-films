package providers

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog"
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/models"
)

func letterboxGetFestivals(config *config.Config, logger *zerolog.Logger, _ *resty.Client) []models.Film {
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
		colly.CacheDir("./cache"),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{DomainGlob: "", Parallelism: 2, RandomDelay: time.Duration(config.DelaySecs) * time.Second})

	var items []models.Film

	c.OnHTML(config.LetterboxProvider.FestivalsSelectorRe, func(h *colly.HTMLElement) {
		link := h.Attr("href")
		re := regexp.MustCompile(`festiville/list`)

		if re.MatchString(link) {
			if config.Debug {
				title := h.Text
				logger.Info().Msgf("Link found: %s, Festival: %s", link, title)
			}
			baseUrl := "https://letterboxd.com" + link
			h.Request.Visit(baseUrl)
		}
	})

	c.OnHTML("li.poster-container", func(h *colly.HTMLElement) {
		id := h.ChildAttr("div.poster", "data-film-id")
		slug := regexp.MustCompile(`-\d{4}(-\d+)?$`).ReplaceAllString(h.ChildAttr("div.poster", "data-film-slug"), "")
		title := strings.Join(strings.Split(slug, "-"), " ")
		idUrlPath := strings.Join(strings.Split(id, ""), "/")
		imageUrl := fmt.Sprintf("https://a.ltrbxd.com/resized/film-poster/%s/%s-%s-0-125-0-187-crop.jpg", idUrlPath, id, slug)

		if config.Debug {
			logger.Info().Msgf("title: %s, slug: %s", title, slug)
		}

		items = append(items, models.Film{
			Provider: "letterbox",
			Id:       id,
			Title:    title,
			ImageUrl: imageUrl,
			Year:     time.Now().Year(),
		})
	})

	if config.Debug {
		c.OnResponse(func(r *colly.Response) {
			logger.Info().Msgf("Response: %s", string(r.Body))
		})
	}

	if config.Debug {
		c.OnRequest(func(r *colly.Request) {
			logger.Info().Msgf("Request to: %s", r.URL.String())
		})
	}

	c.Visit(config.LetterboxProvider.FestivalsUrl)
	c.Wait()
	return items
}

func letterboxGetPopular(config *config.Config, logger *zerolog.Logger, _ *resty.Client) []models.Film {
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
		colly.CacheDir("./cache"),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{DomainGlob: "", Parallelism: 2, RandomDelay: 10 * time.Second})

	var films []models.Film
	c.OnHTML(config.LetterboxProvider.PopularSelectorRe, func(h *colly.HTMLElement) {
		id := h.ChildAttr("div.really-lazy-load", "data-film-id")
		slug := h.ChildAttr("div.really-lazy-load", "data-film-slug")
		idUrlPath := strings.Join(strings.Split(id, ""), "/")
		title := h.ChildAttr("img", "alt")
		imageUrl := fmt.Sprintf("https://a.ltrbxd.com/resized/film-poster/%s/%s-%s-0-140-0-210-crop.jpg", idUrlPath, id, slug)
		films = append(films, models.Film{
			Provider: "letterbox",
			Id:       id,
			Title:    title,
			Year:     time.Now().Year(),
			ImageUrl: imageUrl,
		})
	})

	if config.Debug {
		c.OnResponse(func(r *colly.Response) {
			logger.Info().Msgf("Response: %s", string(r.Body))
		})
	}

	if config.Debug {
		c.OnRequest(func(r *colly.Request) {
			logger.Info().Msgf("Requesting URL: %s", r.URL.String())
		})
	}

	c.Visit(fmt.Sprintf("%s%d%s", config.LetterboxProvider.PopularUrl, time.Now().Year(), "/?esiAllowFilters=true"))
	c.Wait()
	return films
}
