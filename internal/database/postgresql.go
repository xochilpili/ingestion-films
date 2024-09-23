package database

import (
	"database/sql"
	"fmt"
	"strings"

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

func (p *Postgres) InsertFilm(table string, columns []string, item *models.FilmItem) error {

	if ok, err := p.FilmExists(table, item.Title); ok || err != nil {
		if err != nil {
			p.logger.Err(err).Msgf("error while validating film %s", item.Title)
			return err
		}
		return nil
	}

	cols := strings.Join(columns, ",")
	sqlStmt := fmt.Sprintf(`insert into %s (%s) values(%s)`, table, cols, p.computeValues(2))
	_, err := p.db.Exec(sqlStmt, item.Title, item.Year)
	if err != nil {
		p.logger.Err(err).Msgf("error while inserting item: %s", item.Title)
		return err
	}
	return nil
}

func (p *Postgres) computeValues(n int) string {
	nums := make([]string, n)
	for i := 0; i < n; i++ {
		nums[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(nums, ",")
}

func (p *Postgres) FilmExists(table string, title string) (bool, error) {
	if p.db == nil {
		p.logger.Fatal().Msg("database not present!")
	}
	var count int
	sqlStmt := fmt.Sprintf(`select count(1) from %s where title = $1`, table)
	err := p.db.QueryRow(sqlStmt, title).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
