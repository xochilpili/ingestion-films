package main

import (
	"fmt"

	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/logger"
	"github.com/xochilpili/ingestion-films/internal/providers"
)

func main() {

	config := config.New()

	logger := logger.New()

	manager := providers.New(config, logger)

	items := manager.GetFestivals("imdb")
	for _, item := range items {
		fmt.Printf("Title: %s\n", item.Title)
	}
}
