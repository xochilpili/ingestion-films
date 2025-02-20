package providers

/* Festival Struct */
type ImdbFestivalRoot struct {
	Props struct {
		PageProps struct {
			EventName   string `json:"eventName"`
			EditionInfo struct {
				Year               int    `json:"year"`
				ID                 string `json:"id"`
				InstanceWithinYear int    `json:"instanceWithinYear"`
				DateRange          struct {
					StartDate struct {
						DateComponents struct {
							Day      int    `json:"day"`
							Month    int    `json:"month"`
							Year     int    `json:"year"`
							Typename string `json:"__typename"`
						} `json:"dateComponents"`
						DisplayableProperty struct {
							Value struct {
								PlaidHTML string `json:"plaidHtml"`
								Typename  string `json:"__typename"`
							} `json:"value"`
							Typename string `json:"__typename"`
						} `json:"displayableProperty"`
						Typename string `json:"__typename"`
					} `json:"startDate"`
					EndDate struct {
						DateComponents struct {
							Day      int    `json:"day"`
							Month    int    `json:"month"`
							Year     int    `json:"year"`
							Typename string `json:"__typename"`
						} `json:"dateComponents"`
						DisplayableProperty struct {
							Value struct {
								PlaidHTML string `json:"plaidHtml"`
								Typename  string `json:"__typename"`
							} `json:"value"`
							Typename string `json:"__typename"`
						} `json:"displayableProperty"`
						Typename string `json:"__typename"`
					} `json:"endDate"`
					Typename string `json:"__typename"`
				} `json:"dateRange"`
				Event struct {
					Name struct {
						Text     string `json:"text"`
						Typename string `json:"__typename"`
					} `json:"name"`
					Location struct {
						Text     string `json:"text"`
						Typename string `json:"__typename"`
					} `json:"location"`
					Typename string `json:"__typename"`
				} `json:"event"`
				Typename string `json:"__typename"`
			} `json:"editionInfo"`
			Edition struct {
				Awards []struct {
					Text                 string `json:"text"`
					ID                   string `json:"id"`
					NominationCategories struct {
						Edges []struct {
							Node struct {
								Category    interface{} `json:"category"`
								Nominations struct {
									Edges []struct {
										Node struct {
											IsWinner        bool        `json:"isWinner"`
											Notes           interface{} `json:"notes"`
											ForEpisodes     interface{} `json:"forEpisodes"`
											AwardedEntities struct {
												AwardTitles []struct {
													Title struct {
														ID        string `json:"id"`
														TitleText struct {
															Text     string `json:"text"`
															Typename string `json:"__typename"`
														} `json:"titleText"`
														OriginalTitleText struct {
															Text     string `json:"text"`
															Typename string `json:"__typename"`
														} `json:"originalTitleText"`
														PrimaryImage struct {
															URL     string `json:"url"`
															Caption struct {
																PlainText string `json:"plainText"`
																Typename  string `json:"__typename"`
															} `json:"caption"`
															ID       string `json:"id"`
															Width    int    `json:"width"`
															Height   int    `json:"height"`
															Typename string `json:"__typename"`
														} `json:"primaryImage"`
														ReleaseYear *struct{
															Year int `json:"year,omitempty"`
															EndYear int `json:"endYear,omitempty"`
														} `json:"releaseYear,omitempty"`
														CanRate struct {
															IsRatable bool   `json:"isRatable"`
															Typename  string `json:"__typename"`
														} `json:"canRate"`
														RatingsSummary struct {
															AggregateRating float64 `json:"aggregateRating"`
															VoteCount       int     `json:"voteCount"`
															Typename        string  `json:"__typename"`
														} `json:"ratingsSummary"`
														TitleType struct {
															ID       string `json:"id"`
															Typename string `json:"__typename"`
														} `json:"titleType"`
														Typename string `json:"__typename"`
													} `json:"title"`
													Typename string `json:"__typename"`
												} `json:"awardTitles"`
												SecondaryAwardNames []struct {
													Name struct {
														ID       string `json:"id"`
														NameText struct {
															Text     string `json:"text"`
															Typename string `json:"__typename"`
														} `json:"nameText"`
														Typename string `json:"__typename"`
													} `json:"name"`
													Note     interface{} `json:"note"`
													Typename string      `json:"__typename"`
												} `json:"secondaryAwardNames"`
												SecondaryAwardCompanies interface{} `json:"secondaryAwardCompanies"`
												Typename                string      `json:"__typename"`
											} `json:"awardedEntities"`
											Typename string `json:"__typename"`
										} `json:"node"`
										Typename string `json:"__typename"`
									} `json:"edges"`
									PageInfo struct {
										HasNextPage bool   `json:"hasNextPage"`
										EndCursor   string `json:"endCursor"`
										Typename    string `json:"__typename"`
									} `json:"pageInfo"`
									Typename string `json:"__typename"`
								} `json:"nominations"`
								Typename string `json:"__typename"`
							} `json:"node"`
							Typename string `json:"__typename"`
						} `json:"edges"`
						Typename string `json:"__typename"`
					} `json:"nominationCategories"`
					Typename string `json:"__typename"`
				} `json:"awards"`
				Typename string `json:"__typename"`
			} `json:"edition"`
			HistoryEventEditions []struct {
				Year               int    `json:"year"`
				ID                 string `json:"id"`
				InstanceWithinYear int    `json:"instanceWithinYear"`
				Typename           string `json:"__typename"`
			} `json:"historyEventEditions"`
			TopEventsData map[string]interface {
			} `json:"topEventsData,omitempty"`
			EventLinks []struct {
				Category struct {
					CategoryID string `json:"categoryId"`
					ID         string `json:"id"`
					Text       string `json:"text"`
					Typename   string `json:"__typename"`
				} `json:"category"`
				URL      string `json:"url"`
				Typename string `json:"__typename"`
			} `json:"eventLinks"`
			EventAkas []struct {
				Name struct {
					Text     string `json:"text"`
					Typename string `json:"__typename"`
				} `json:"name"`
				Typename string `json:"__typename"`
			} `json:"eventAkas"`
			EventDateRange map[string]interface {
				} `json:"eventDateRange,omitempty"`
			RequestContext map[string]interface{
				} `json:"requestContext,omitempty"`
			CmsContext map[string]interface{} `json:"cmsContext,omitempty"`
			TranslationContext map[string]interface{} `json:"translationContext,omitempty"`
			UrqlState  interface{} `json:"urqlState"`
			FetchState interface{} `json:"fetchState"`
		} `json:"pageProps"`
		NSSP bool `json:"__N_SSP"`
	} `json:"props"`
	Page  string `json:"page"`
	Query map[string]interface {
	} `json:"query,omitempty"`
	BuildID       string `json:"buildId"`
	AssetPrefix   string `json:"assetPrefix"`
	RuntimeConfig map[string]interface {
	} `json:"runtimeConfig,omitempty"`
	IsFallback    bool          `json:"isFallback"`
	Gssp          bool          `json:"gssp"`
	Locale        string        `json:"locale"`
	Locales       []string      `json:"locales"`
	DefaultLocale string        `json:"defaultLocale"`
	ScriptLoader  []interface{} `json:"scriptLoader"`
}

/* Popular Struct */
type ImdbPopularGenres struct{
	Genre struct{
		Text string `json:"text"`
		Typename string `json:"__typename"`
	} `json:"genre"`
	Typename string `json:"__typename"`
}

type ImdbPopularRoot struct{
		Props struct {
			PageProps struct {
				InitialRefinerQueryInfo map[string]interface{} `json:"initialRefinerQueryInfo"`
				PageData struct {
					ChartTitles struct {
						Edges []struct {
							CurrentRank int `json:"currentRank"`
							Node        struct {
								ID        string `json:"id"`
								TitleText struct {
									Text     string `json:"text"`
									Typename string `json:"__typename"`
								} `json:"titleText"`
								TitleType struct {
									ID                  string `json:"id"`
									Text                string `json:"text"`
									CanHaveEpisodes     bool   `json:"canHaveEpisodes"`
									DisplayableProperty struct {
										Value struct {
											PlainText string `json:"plainText"`
											Typename  string `json:"__typename"`
										} `json:"value"`
										Typename string `json:"__typename"`
									} `json:"displayableProperty"`
									Typename string `json:"__typename"`
								} `json:"titleType"`
								OriginalTitleText struct {
									Text     string `json:"text"`
									Typename string `json:"__typename"`
								} `json:"originalTitleText"`
								PrimaryImage struct {
									ID      string `json:"id"`
									Width   int    `json:"width"`
									Height  int    `json:"height"`
									URL     string `json:"url"`
									Caption struct {
										PlainText string `json:"plainText"`
										Typename  string `json:"__typename"`
									} `json:"caption"`
									Typename string `json:"__typename"`
								} `json:"primaryImage"`
								ReleaseYear struct {
									Year     int    `json:"year"`
									EndYear  any    `json:"endYear"`
									Typename string `json:"__typename"`
								} `json:"releaseYear"`
								RatingsSummary struct {
									AggregateRating float64 `json:"aggregateRating"`
									VoteCount       int     `json:"voteCount"`
									Typename        string  `json:"__typename"`
								} `json:"ratingsSummary"`
								Runtime struct {
									Seconds  int    `json:"seconds"`
									Typename string `json:"__typename"`
								} `json:"runtime"`
								Certificate struct {
									Rating   string `json:"rating"`
									Typename string `json:"__typename"`
								} `json:"certificate"`
								CanRate struct {
									IsRatable bool   `json:"isRatable"`
									Typename  string `json:"__typename"`
								} `json:"canRate"`
								TitleGenres struct {
									Genres []ImdbPopularGenres `json:"genres"`
									Typename string `json:"__typename"`
								} `json:"titleGenres"`
								CanHaveEpisodes bool `json:"canHaveEpisodes"`
								Plot            struct {
									PlotText struct {
										PlainText string `json:"plainText"`
										Typename  string `json:"__typename"`
									} `json:"plotText"`
									Typename string `json:"__typename"`
								} `json:"plot"`
								LatestTrailer struct {
									ID       string `json:"id"`
									Typename string `json:"__typename"`
								} `json:"latestTrailer"`
								Series      any `json:"series"`
								ReleaseDate struct {
									Day      int    `json:"day"`
									Month    int    `json:"month"`
									Year     int    `json:"year"`
									Typename string `json:"__typename"`
								} `json:"releaseDate"`
								ProductionStatus struct {
									CurrentProductionStage struct {
										ID       string `json:"id"`
										Text     string `json:"text"`
										Typename string `json:"__typename"`
									} `json:"currentProductionStage"`
									Typename string `json:"__typename"`
								} `json:"productionStatus"`
								MeterRanking struct {
									CurrentRank int `json:"currentRank"`
									RankChange  struct {
										ChangeDirection string `json:"changeDirection"`
										Difference      int    `json:"difference"`
										Typename        string `json:"__typename"`
									} `json:"rankChange"`
									Typename string `json:"__typename"`
								} `json:"meterRanking"`
								PrincipalCredits []map[string]interface{} `json:"principalCredits"`
								Typename string `json:"__typename"`
							} `json:"node"`
							Typename string `json:"__typename"`
						} `json:"edges"`
						Genres []struct {
							FilterID string `json:"filterId"`
							Text     string `json:"text"`
							Total    int    `json:"total"`
							Typename string `json:"__typename"`
						} `json:"genres"`
						Keywords []struct {
							FilterID string `json:"filterId"`
							Text     string `json:"text"`
							Total    int    `json:"total"`
							Typename string `json:"__typename"`
						} `json:"keywords"`
						WatchOptions []struct {
							FilterID string `json:"filterId"`
							Text     string `json:"text"`
							Total    int    `json:"total"`
							Typename string `json:"__typename"`
						} `json:"watchOptions"`
						Typename string `json:"__typename"`
					} `json:"chartTitles"`
				} `json:"pageData"`
				RequestContext map[string]interface {} `json:"requestContext"`
				CmsContext map[string]interface{} `json:"cmsContext"`
				TranslationContext map[string]interface{} `json:"translationContext"`
				UrqlState  any `json:"urqlState"`
				FetchState any `json:"fetchState"`
			} `json:"pageProps"`
			NSsp bool `json:"__N_SSP"`
		} `json:"props"`
		Page  string `json:"page"`
		Query map[string]interface{} `json:"query"`
		BuildID       string `json:"buildId"`
		AssetPrefix   string `json:"assetPrefix"`
		RuntimeConfig map[string]interface{
		} `json:"runtimeConfig"`
		IsFallback    bool     `json:"isFallback"`
		Gssp          bool     `json:"gssp"`
		Locale        string   `json:"locale"`
		Locales       []string `json:"locales"`
		DefaultLocale string   `json:"defaultLocale"`
		ScriptLoader  []any    `json:"scriptLoader"`
}

