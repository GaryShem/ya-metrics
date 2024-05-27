package handlers

import (
	"github.com/GaryShem/ya-metrics.git/internal/server/storage/repository"
)

type RepoHandler struct {
	repo repository.Repository
}

func NewHandler(repo repository.Repository) *RepoHandler {
	return &RepoHandler{repo: repo}
}
