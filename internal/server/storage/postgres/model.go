package postgres

import (
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
)

var ErrSQLConnectionFailed = errors.New("SQL connection failed")
var ErrSQLNotImplemented = errors.New("SQL not implemented yet")

type SQLStorage struct {
	ConnString string
	db         *sql.DB
}

func NewSQLStorage(conn string) *SQLStorage {
	return &SQLStorage{
		ConnString: conn,
	}
}

func (s *SQLStorage) Init() error {
	db, err := sql.Open("pgx", s.ConnString)
	if err != nil {
		return err
	}
	// create tables
	createGaugeTableSQL := `CREATE TABLE IF NOT EXISTS gauges ( 
    id varchar(45) NOT NULL,
    val double precision NOT NULL,
	PRIMARY KEY (id));`
	if _, err = db.Exec(createGaugeTableSQL); err != nil {
		return err
	}
	createCounterTableSQL := `CREATE TABLE IF NOT EXISTS counters ( 
    id varchar(45) NOT NULL,
    val bigint NOT NULL,
	PRIMARY KEY (id));`
	if _, err = db.Exec(createCounterTableSQL); err != nil {
		return err
	}
	s.db = db
	return nil
}

var _ repository.Repository = NewSQLStorage("")
