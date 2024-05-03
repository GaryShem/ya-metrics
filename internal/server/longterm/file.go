package longterm

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/GaryShem/ya-metrics.git/internal/server/storage/memorystorage"
	"github.com/GaryShem/ya-metrics.git/internal/shared/logging"
)

type FileSaver struct {
	timestamp time.Time
	filename  string
	MS        *memorystorage.MemStorage
}

func NewFileSaver(filename string, ms *memorystorage.MemStorage) *FileSaver {
	return &FileSaver{
		filename: filename,
		MS:       ms,
	}
}

func (fs *FileSaver) SaveMetricsFile(interval time.Duration) error {
	if len(fs.filename) < 1 {
		return nil
	}
	if interval == 0 {
		interval = time.Millisecond * 100
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if !fs.timestamp.Before(fs.MS.LastChangeTime) {
			continue
		}
		j, err := json.Marshal(fs.MS)
		if err != nil {
			return fmt.Errorf("save metrics file: %w", err)
		}
		if err = os.WriteFile(fs.filename, j, 0644); err != nil {
			return fmt.Errorf("save metrics file: %w", err)
		}
	}
	return nil
}

func (fs *FileSaver) LoadMetricsFile() error {
	ms := &memorystorage.MemStorage{}
	j, err := os.ReadFile(fs.filename)
	if err != nil {
		logging.Log.Infoln("unable to load metrics file:", err)
		fs.MS = memorystorage.NewMemStorage()
		return nil
	}
	if err := json.Unmarshal(j, ms); err != nil {
		return fmt.Errorf("load metrics file: %w", err)
	}
	fs.timestamp = ms.LastChangeTime
	fs.MS = ms
	return nil
}
