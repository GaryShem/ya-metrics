package handlers

import (
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage"
)

type RepoHandler struct {
	repo storage.Repository
}

func NewHandler(repo storage.Repository) *RepoHandler {
	return &RepoHandler{repo: repo}
}
