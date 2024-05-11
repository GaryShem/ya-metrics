package postgres

type PostgreSQLStorage struct {
	ConnString string
}

func NewPostgreSQLStorage(conn string) *PostgreSQLStorage {
	return &PostgreSQLStorage{
		ConnString: conn,
	}
}
