package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *SQLStorage) Ping() error {
	db, err := sql.Open("pgx", s.ConnString)
	if err != nil {
		return ErrSQLConnectionFailed
	}
	defer db.Close()
	return db.Ping()
}
