package main

import (
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/logger"
	"github.com/xochilpili/ingestion-films/internal/providers"
)

func main() {

	config := config.New()

	logger := logger.New()

	letterbox := providers.NewLetterBox(config, logger)

	films := letterbox.GetFestivals()
	logger.Info().Msgf("Films: %v", films)
}
