package postgres

import (
	"database/sql"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (pss *PostgreSQLStorage) TestConnection(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("pgx", pss.ConnString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	w.WriteHeader(http.StatusOK)
}
