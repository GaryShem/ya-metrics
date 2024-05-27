package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (pss *PostgreSQLStorage) Ping() error {
	db, err := sql.Open("pgx", pss.ConnString)
	if err != nil {
		return ErrSQLConnectionFailed
	}
	defer db.Close()
	return db.Ping()
}
