package models

type Film struct {
	Id          string
	Title       string
	Description string
	ImageUrl    string
	Year        int
	Genre       []string
}

type FilmItem struct {
	Provider string
	Title    string
	Year     int
}
