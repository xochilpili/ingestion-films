package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xochilpili/ingestion-films/internal/models"
)

func (w *WebServer) GetPopularByProvider(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, &gin.H{"message": "error", "errors": "bad request"})
		return
	}

	w.logger.Info().Msgf("request Popular films by %s provider", provider)
	films := w.manager.GetPopular(provider)
	w.logger.Info().Msgf("resolved %d popular films", len(films))
	c.JSON(http.StatusOK, &gin.H{"message": "ok", "total": len(films), "films": films})
}

func (w *WebServer) GetAllPopular(c *gin.Context) {
	providers := []string{"yts", "imdb", "letterbox"}

	w.logger.Info().Msgf("request Popular films by all providers")
	var films []models.Film

	for _, provider := range providers {
		items := w.manager.GetPopular(provider)
		w.logger.Info().Msgf("returned %d films by %s provider", len(items), provider)
		films = append(films, items...)
	}
	w.logger.Info().Msgf("resolved %d popular films", len(films))
	c.JSON(http.StatusOK, &gin.H{"message": "ok", "total": len(films), "films": films})
}
