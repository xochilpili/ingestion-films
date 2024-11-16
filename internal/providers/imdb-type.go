package providers

/* Festival Struct */
type ImdbFestivalRootObject struct {
	NomineesWidgetModel ImdbFestivalAwardNominations `json:"nomineesWidgetModel"`
}

type ImdbFestivalAwardNominations struct {
	RefMarker               string                     `json:"refMarker"`
	EventEditionSummary     ImdbFestivalEditionSummary `json:"eventEditionSummary"`
	AlwaysDisplayAwardNames bool                       `json:"alwaysDisplayAwardNames"`
	ShouldDisplayNotes      bool                       `json:"shouldDisplayNotes"`
	Mobile                  bool                       `json:"mobile"`
}

type ImdbFestivalEditionSummary struct {
	Awards     []ImdbFestivalAwards `json:"awards,omitempty"`
	EventId    string               `json:"eventId"`
	Occurrence int                  `json:"occurrence"`
	RequestKey struct {
		Year              int    `json:"year"`
		Occurrence        int    `json:"occurrence"`
		PreferredLanguage string `json:"preferredLanguage"`
		PreferredRegion   string `json:"preferredRegion"`
		UserRegion        string `json:"userRegion"`
		EventId           string `json:"eventId"`
	} `json:"requestKey"`
	EventName      string `json:"eventName"`
	EventEditionId string `json:"eventEditionId"`
	Year           int    `json:"year"`
}

type ImdbFestivalAwards struct {
	Id         string                   `json:"id"`
	AwardName  string                   `json:"awardName"`
	Trivia     []string                 `json:"trivia,omitempty"`
	Categories []ImdbFestivalCategories `json:"categories"`
}

type ImdbFestivalCategories struct {
	CategoryName      string                    `json:"categoryName"`
	Nominations       []ImdbFestivalNominations `json:"nominations"`
	IsPrimaryCategory bool                      `json:"isPrimaryCategory"`
}

type ImdbFestivalNominations struct {
	PrimaryNominees     []ImdbFestivalNominees `json:"primaryNominees"`
	SecondayNominees    []ImdbFestivalNominees `json:"secondaryNominees"`
	CategoryName        string                 `json:"categoryName,omitempty"`
	WinAnnouncementTime int                    `json:"winAnnouncementTime,omitempty"`
	Notes               string                 `json:"notes,omitempty"`
	CharactedNames      []string               `json:"characterNames,omitempty"`
	SongNames           []string               `json:"songNames,omitempty"`
	EpisodeNames        []string               `json:"episodeNames,omitempty"`
	AwardName           string                 `json:"awardName"`
	AwardNominationId   string                 `json:"awardNominationId"`
	IsWinner            bool                   `json:"isWinner"`
}

type ImdbFestivalNominees struct {
	Name         string `json:"name"`
	Note         string `json:"note,omitempty"`
	ImageUrl     string `json:"imageUrl"`
	ImageHeight  int    `json:"imageHeight"`
	ImageWidth   int    `json:"imageWidth"`
	OriginalName string `json:"originalName,omitempty"`
	Const        string `json:"const"`
}

/* Popular Struct */

type ImdbPopularRootObject struct {
	Type            string               `json:"@type"`
	ItemListElement []ImdbPopularElement `json:"itemListElement"`
	Context         string               `json:"@context"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
}

type ImdbPopularElement struct {
	Type string          `json:"@type"`
	Item ImdbPopularItem `json:"item"`
}

type ImdbPopularItem struct {
	Type            string                 `json:"@type"`
	Url             string                 `json:"url"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Image           string                 `json:"image"`
	AggregateRating ImdbPopularItemRaiting `json:"aggregateRating"`
	ContentRaiting  string                 `json:"contentRaiting"`
	Genre           string                 `json:"genre"`
	Duration        string                 `json:"duration"`
	Year            string                 `json:"year,omitempty"`
}

type ImdbPopularItemRaiting struct {
	Type        string  `json:"@type"`
	BestRating  int     `json:"bestRating"`
	WorstRating int     `json:"worstRating"`
	RatingValue float64 `json:"ratingValue"`
	RatingCount int     `json:"ratingCount"`
}
