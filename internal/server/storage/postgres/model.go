package postgres

import "errors"

var SQLNotInitialized = errors.New("SQL not initialized")
var SQLConnectionFailed = errors.New("SQL connection failed")

var SQLStorage *PostgreSQLStorage

type PostgreSQLStorage struct {
	ConnString string
}

func NewPostgreSQLStorage(conn string) *PostgreSQLStorage {
	return &PostgreSQLStorage{
		ConnString: conn,
	}
}
