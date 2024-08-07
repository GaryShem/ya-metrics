package handlers

import "net/http"

// Ping - heartbeat method for the RepoHandler.
func (h *RepoHandler) Ping(w http.ResponseWriter, _ *http.Request) {
	err := h.repo.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
