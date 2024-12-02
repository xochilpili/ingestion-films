package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Database struct {
	Host     string `required:"true"`
	Name     string `default:"ingestion_films"`
	User     string `required:"true"`
	Password string `required:"true"`
}

type Tmdb struct {
	Url    string `required:"true" split_words:"true"`
	ApiKey string `required:"true" split_words:"true"`
}

type Plex struct {
	ApiUrl string `required:"true" split_words:"true"`
}

type Config struct {
	HOST          string   `default:"0.0.0.0" required:"true"`
	PORT          string   `default:"4000" required:"true"`
	ServiceName   string   `default:"ingestion-films" split_words:"true"`
	Debug         bool     `default:"false"`
	Database      Database `split_words:"true" required:"true"`
	Tmdb          Tmdb     `required:"true"`
	ExcludeGenres []string `split_words:"true"`
	Plex          Plex     `split_words:"true"`
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
