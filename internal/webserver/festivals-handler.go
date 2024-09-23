package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (w *WebServer) GetFestivalsByProvider(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" || provider == "yts" {
		c.JSON(http.StatusBadRequest, &gin.H{"message": "error", "error": "bad request"})
		return
	}

	w.logger.Info().Msgf("requesting Festivals from %s provider", provider)
	films := w.manager.GetFestivals(provider)
	w.logger.Info().Msgf("resolved %d films", len(films))
	c.JSON(http.StatusOK, &gin.H{"message": "ok", "total": len(films), "films": films})
}

func (w *WebServer) GetAllFestivals(c *gin.Context) {

	w.logger.Info().Msgf("requesting festivals by all providers")

	films := w.manager.GetFestivals("")

	w.logger.Info().Msgf("resolved %d films", len(films))
	c.JSON(http.StatusOK, &gin.H{"message": "ok", "total": len(films), "films": films})
}

func (w *WebServer) SyncFestivals(c *gin.Context) {
	provider := c.Param("provider")

	w.logger.Info().Msg("request syncing festivals")
	go func() {
		err := w.manager.SyncFestivals(provider)
		if err != nil {
			c.JSON(http.StatusBadRequest, &gin.H{"message": "error", "error": err.Error()})
			return
		}
	}()

	c.JSON(http.StatusOK, &gin.H{"message": "ok"})
}
