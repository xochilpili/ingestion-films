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
	Title string
	Year  int
}
