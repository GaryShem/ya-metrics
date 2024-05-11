package memorystorage

import (
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/postgres"
)

func (ms *MemStorage) Ping() error {
	if postgres.SQLStorage == nil {
		return postgres.SQLNotInitialized
	}
	return postgres.SQLStorage.Ping()
}
