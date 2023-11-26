package handlers

import (
	"net/http"
	"pass-it/internal/server/templates"
)

type IndexRequestHandler struct {
}

func NewIndexRequestHandler() http.Handler {
	return &IndexRequestHandler{}
}

func (h *IndexRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	templates.Index().Render(r.Context(), w)
}
