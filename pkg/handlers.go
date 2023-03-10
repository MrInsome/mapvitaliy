package pkg

import "net/http"

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
