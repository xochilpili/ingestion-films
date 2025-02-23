package providers

type YtsPopularRootObject struct {
	Status        string         `json:"status"`
	StatusMessage string         `json:"status_message"`
	Data          YtsPopularData `json:"data"`
	Meta          struct {
		ServerTime     int    `json:"server_time"`
		ServerTimezone string `json:"server_timezone"`
		ApiVersion     int    `json:"api_version"`
		Executiontime  string `json:"execution_time"`
	} `json:"@meta"`
}

type YtsPopularData struct {
	MovieCount int       `json:"movie_count"`
	Limit      int       `json:"limit"`
	PageNumber int       `json:"page_number"`
	Movies     []YtsFilm `json:"movies"`
}

type YtsFilm struct {
	Id                      int          `json:"id"`
	Url                     string       `json:"url"`
	ImdbCode                string       `json:"imdb_code"`
	Title                   string       `json:"title"`
	TitleEnglish            string       `json:"title_english"`
	TitleLong               string       `json:"title_long"`
	Slug                    string       `json:"slug"`
	Year                    int          `json:"year"`
	Rating                  float64      `json:"rating"`
	Runtime                 int          `json:"runtime"`
	Genres                  []string     `json:"genres"`
	Summary                 string       `json:"summary,omitempty"`
	DescriptionFull         string       `json:"description_full,omitempty"`
	Synopsis                string       `json:"synopsis,omitempty"`
	YtTrailerCode           string       `json:"yt_trailer_code,omitempty"`
	Language                string       `json:"language"`
	MpaRating               string       `json:"mpa_rating,omitempty"`
	BackgroundImage         string       `json:"background_image"`
	BackgroundImageOriginal string       `json:"background_image_original"`
	SmallCoverImage         string       `json:"small_cover_image"`
	MediumCoverImage        string       `json:"medium_cover_image"`
	LargeCoverImage         string       `json:"large_cover_image"`
	State                   string       `json:"state"`
	Torrents                []YtsTorrent `json:"torrents"`
	DateUploaded            string       `json:"date_uploaded"`
	DateUploadedUnix        int          `json:"date_uploaded_unix"`
}

type YtsTorrent struct {
	Url              string `json:"url"`
	Hash             string `json:"hash"`
	Quality          string `json:"quality"`
	Type             string `json:"type"`
	IsRepack         string `json:"is_repack"`
	VideoCodec       string `json:"video_codec"`
	BitDepth         string `json:"bit_depth"`
	AudioChannels    string `json:"audio_channels"`
	Seeds            int    `json:"seeds"`
	Peers            int    `json:"peers"`
	Size             string `json:"size"`
	SizeBytes        int    `json:"size_bytes"`
	DateUploaded     string `json:"date_uploaded"`
	DateUploadedUnix int    `json:"date_uploaded_unix"`
}
