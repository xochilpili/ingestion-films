package models

type FestivalFilm struct {
	Id            string
	Title         string
	OriginalTitle string
	Year          int
	ImageUrl      string
	Festival      string
}

type PopularFilm struct {
	Id          string
	Title       string
	Description string
	ImageUrl    string
	Genre       []string
}
