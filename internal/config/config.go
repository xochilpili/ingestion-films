package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type ImdbProvider struct {
	HttpPrefix        string            `default:"https://" split_words:"true"`
	Festivals         map[string]string `default:"cannes:www.imdb.com/event/ev0000147/,tiff:www.imdb.com/event/ev0000659/,venecia:www.imdb.com/event/ev0000681/,oscar:www.imdb.com/event/ev0000003/,berlinale:www.imdb.com/event/ev0000091/"`
	PopularUrl        string            `default:"https://www.imdb.com/chart/moviemeter/?ref_=nv_mv_mpm" split_words:"true"`
	PopularSelectorRe string            `default:"itemListElement" split_words:"true"`
	DelaySecs         int               `default:"10" split_words:"true"`
}

type Config struct {
	HOST         string       `default:"0.0.0.0"`
	PORT         string       `default:"4000"`
	ServiceName  string       `default:"ingestion-films" split_words:"true"`
	Debug        bool         `default:"false"`
	ImdbProvider ImdbProvider `split_words:"true"`
}

var EnvPrefix = "IFS"

func New() *Config {
	godotenv.Load()
	cfg, err := Get()
	if err != nil {
		panic(fmt.Errorf("invalid value(s) retrieved from environment: %w", err))
	}
	return cfg
}

func Get() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process(EnvPrefix, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
