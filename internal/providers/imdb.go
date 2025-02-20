package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog"
	"github.com/xochilpili/ingestion-films/internal/models"
)

func imdbGetFestivals(config *ProviderConfig, logger *zerolog.Logger, _ *resty.Client) []models.Film {
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
		colly.CacheDir("./cache"),
		colly.UserAgent(config.UserAgent),
	)

	c.Limit(&colly.LimitRule{DomainGlob: "", Parallelism: 2, RandomDelay: time.Duration(config.DelaySecs) * time.Second})

	festivals := config.Festivals

	/*
		festivals := map[string]string{
			"cannes":    "https://www.imdb.com/event/ev0000147/",
			"tiff":      "https://www.imdb.com/event/ev0000659/",
			"venecia":   "https://www.imdb.com/event/ev0000681/",
			"oscar":     "https://www.imdb.com/event/ev0000003/",
			"berlinale": "https://www.imdb.com/event/ev0000091/",
		}
	*/

	var model ImdbFestivalRoot
	var films []models.Film

	for k, url := range festivals {
		
		c.OnHTML("script[id=__NEXT_DATA__]", func(h *colly.HTMLElement) {
			scriptContent := h.Text

			if config.Debug {
				fmt.Printf("%v", scriptContent)
			}

			err := json.Unmarshal([]byte(scriptContent), &model)
				if err != nil {
					log.Fatalf("error while unmarshal %v\n", err)
				}
			films = append(films, translate2FestivalModel(&model)...)
		})

		if config.Debug {
			c.OnRequest(func(r *colly.Request) {
				fmt.Printf("Requesting festival: %s, %s\n", k, r.URL.String())
			})
		}

		logger.Info().Msgf("Visiting Festival: %s, url: %s", k, url)

		c.Visit(url)
		c.Wait()
	}

	return films
}

func imdbGetPopular(config *ProviderConfig, logger *zerolog.Logger, _ *resty.Client) []models.Film {
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
		colly.CacheDir("./cache"),
		colly.UserAgent(config.UserAgent),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1, RandomDelay: time.Duration(config.DelaySecs) * time.Second})

	imdbPopularUrl := config.PopularUrl

	var model ImdbPopularRoot

	c.OnHTML("script[id=__NEXT_DATA__]", func(h *colly.HTMLElement) {
		scriptContent := h.Text

		if config.Debug {
			fmt.Printf("%s\n", scriptContent)
		}

		err := json.Unmarshal([]byte(scriptContent), &model)
		if err != nil {
			log.Fatalf("error while unmarshal %v\n", err)
		}
		/* decoder := json.NewDecoder(strings.NewReader(scriptContent))
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&model)
		if err != nil{
			panic(fmt.Errorf("error unmarshal %v", err))
		} */

	})

	if config.Debug {
		c.OnResponse(func(r *colly.Response) {
			fmt.Printf("Response:\n")
			fmt.Println(string(r.Body))
		})
	}
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

/*
* IMDB changed their schema :(
*/
func translate2FestivalModel(imdbObject *ImdbFestivalRoot) []models.Film {
	var films []models.Film
	for _, award := range(imdbObject.Props.PageProps.Edition.Awards){
		for _, node := range(award.NominationCategories.Edges){
			for _, edge := range(node.Node.Nominations.Edges){
				for _, item := range(edge.Node.AwardedEntities.AwardTitles){
					films = append(films, models.Film{
						Provider: "imdb",
						Id: item.Title.ID,
						Title: sanitizeTitle(item.Title.TitleText.Text),
						Year: 0,
						ImageUrl: item.Title.PrimaryImage.URL,
					})
				}
			}
		}
	}
	return films
}

func translate2PopularModel(imdbObject *ImdbPopularRoot) []models.Film {
	var films []models.Film
	for _, item := range imdbObject.Props.PageProps.PageData.ChartTitles.Edges{
		films = append(films, models.Film{
			Provider: "imdb",
			Id: item.Node.ID,
			Title: sanitizeTitle(item.Node.OriginalTitleText.Text),
			Description: "",
			ImageUrl: item.Node.PrimaryImage.URL,
			Year: item.Node.ReleaseYear.Year,
			Genre: popularGenres(item.Node.TitleGenres.Genres),
		})	
	}
	return films
}

func popularGenres(arr []ImdbPopularGenres) []string{
	genres := make([]string, len(arr))
	for i, genre := range(arr){
		genres[i] = genre.Genre.Text
	}
	return genres
}

func sanitizeTitle(title string) string{
	str := strings.ReplaceAll(title, "&amp;", "&")
	str = strings.ReplaceAll(str, "&apos;", "'")
	return str
}

func imdbPostProcess(config *ProviderConfig, pgService PgService, apiService ApiService, items *[]models.Film) (*[]models.Film, error) {
	if !config.RequireTmdb {
		return items, nil
	}
	for i := range *items {
		item := &(*items)[i]
		if config.Debug {
			fmt.Printf("getting item details: %s\n", item.Title)
		}
		filmdetails, err := apiService.GetMovieDetails(item.Title)
		if err != nil {
			return nil, err
		}
		for _, film := range filmdetails {
			if strings.EqualFold(film.Title, item.Title) {
				yearStr := strings.Split(film.ReleaseDate, "-")
				year, _ := strconv.Atoi(yearStr[0])
				if len(film.GenreIds) == 0 {
					item.Genre = nil
					item.Year = year
					continue
				}
				genres, err := pgService.GetGenres(film.GenreIds)
				if err != nil {
					continue
				}
				item.Description = film.Overview
				item.Year = year
				item.Genre = genres
			}
		}
	}
	return items, nil
}
