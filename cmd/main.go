package main

import (
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/logger"
	"github.com/xochilpili/ingestion-films/internal/providers"
)

func main() {

	config := config.New()

	logger := logger.New()

	imdbProvider := providers.NewImdb(config, logger)

	films := imdbProvider.GetFestivals()
	logger.Info().Msgf("films: %v\n", films)
}
