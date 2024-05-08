package handlers

import (
	"github.com/GaryShem/ya-metrics.git/internal/shared/storage/models"
)

type RepoHandler struct {
	repo models.Repository
}

func NewHandler(repo models.Repository) *RepoHandler {
	return &RepoHandler{repo: repo}
}
