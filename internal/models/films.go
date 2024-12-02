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

type TmdbItem struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIds         []int   `json:"genre_ids"`
	Id               int     `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"oroginal_title"`
	Overview         string  `json:"overview"`
	Popularity       float64 `json:"popularity"`
	PosterPath       string  `json:"poster_path"`
	ReleaseDate      string  `json:"release_date"`
	Title            string  `json:"title"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
}

type TmdbResponse struct {
	Page         int        `json:"page"`
	Results      []TmdbItem `json:"results"`
	TotalPages   int        `json:"total_pages"`
	TotalResults int        `json:"total_results"`
}

type PlexMediaItem struct {
	Title         string `json:"title"`
	Type          string `json:"type"`
	HubIdentifier string `json:"hubIdentifier"`
	Metadata      struct {
		Key                 string `json:"key"`
		Type                string `json:"type"`
		LibrarySectionID    int    `json:"librarySectionID"`
		LibrarySectionTitle string `json:"librarySectionTitle"`
		Title               string `json:"title"`
		OriginalTitle       string `json:"originalTitle"`
		Year                int    `json:"year"`
		Location            []struct {
			Path string `json:"path"`
		} `json:"Location"`
	} `json:"Metadata"`
}

type PlexMediaSearch struct {
	MediaContainer struct {
		Size         int             `json:"size,omitempty"`
		SearchResult []PlexMediaItem `json:"SearchResult"`
	} `json:"MediaContainer,omitempty"`
}
