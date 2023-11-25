package handlers

import (
	"net/http"
	"pass-it/server/templates"
)

type IndexRequestHandler struct {
}

func NewIndexRequestHandler() *IndexRequestHandler {
	return &IndexRequestHandler{}
}

func (h *IndexRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	templates.Index().Render(r.Context(), w)
}
