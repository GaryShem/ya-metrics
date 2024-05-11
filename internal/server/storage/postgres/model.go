package postgres

import (
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
)

var ErrSQLConnectionFailed = errors.New("SQL connection failed")
var ErrSQLNotImplemented = errors.New("SQL not implemented yet")

type SqlStorage struct {
	ConnString string
	db         *sql.DB
}

func NewSqlStorage(conn string) *SqlStorage {
	return &SqlStorage{
		ConnString: conn,
	}
}

func (s *SqlStorage) Init() error {
	db, err := sql.Open("pgx", s.ConnString)
	if err != nil {
		return err
	}
	// create tables
	createGaugeTableSql := `CREATE TABLE IF NOT EXISTS gauges ( 
    id varchar(45) NOT NULL,
    val double precision NOT NULL,
	PRIMARY KEY (id));`
	if _, err = db.Exec(createGaugeTableSql); err != nil {
		return err
	}
	createCounterTableSql := `CREATE TABLE IF NOT EXISTS counters ( 
    id varchar(45) NOT NULL,
    val bigint NOT NULL,
	PRIMARY KEY (id));`
	if _, err = db.Exec(createCounterTableSql); err != nil {
		return err
	}
	s.db = db
	return nil
}

var _ repository.Repository = NewSqlStorage("")
