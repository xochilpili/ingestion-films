package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/xochilpili/ingestion-films/internal/config"
	"github.com/xochilpili/ingestion-films/internal/models"
)

type Postgres struct {
	config *config.Config
	logger *zerolog.Logger
	db     *sql.DB
}

func New(config *config.Config, logger *zerolog.Logger) *Postgres {
	return &Postgres{
		config: config,
		logger: logger,
	}
}

func (p *Postgres) Connect() error {
	psqlConn := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", p.config.Database.Host, p.config.Database.User, p.config.Database.Password, p.config.Database.Name)
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return err
	}
	//defer db.Close()
	p.db = db
	return nil
}

func (p *Postgres) Ping() error {
	return p.db.Ping()
}

func (p *Postgres) Close() error {
	err := p.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) InsertFestivalFilm(film *models.FestivalItem) error {
	//defer p.db.Close()
	if p.db == nil {
		p.logger.Fatal().Msg("database not present!")
	}

	if ok, err := p.FilmExists(film.Title); ok || err != nil {
		if err != nil {
			p.logger.Err(err).Msgf("error while validating film %s", film.Title)
			return err
		}
		return nil
	}

	sqlStmt := `insert into films_festivals (title,year) values($1, $2)`
	_, err := p.db.Exec(sqlStmt, film.Title, film.Year)
	if err != nil {
		p.logger.Err(err).Msgf("error while inserting film %s", film.Title)
		return err
	}

	return nil
}

func (p *Postgres) FilmExists(title string) (bool, error) {
	if p.db == nil {
		p.logger.Fatal().Msg("database not present!")
	}
	var count int
	sqlStmt := `select count(1) from films_festivals where title = $1`
	err := p.db.QueryRow(sqlStmt, title).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
