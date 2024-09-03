package http

import (
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
)

// RepoHandler - handler structure for the storage.
type RepoHandler struct {
	repo repository.Repository
}

// NewHandler - RepoHandler constructor function.
func NewHandler(repo repository.Repository) *RepoHandler {
	return &RepoHandler{repo: repo}
}
