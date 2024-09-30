package models

type Film struct {
	Provider    string
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
