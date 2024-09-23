package webserver

func (w *WebServer) loadRoutes() {

	api := w.ginger.Group("/")
	api.GET("/ping", w.PingHandler)
	festivals := w.ginger.Group("/festivals")
	{
		festivals.GET("/:provider", w.GetFestivalsByProvider)
		festivals.GET("/all", w.GetAllFestivals)
		festivals.GET("/:provider/sync", w.SyncFestivals)
		festivals.GET("/all/sync", w.SyncFestivals)
	}
	popular := w.ginger.Group("/popular")
	{
		popular.GET("/:provider", w.GetPopularByProvider)
		popular.GET("/all", w.GetAllPopular)
	}
}
