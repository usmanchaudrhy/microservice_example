package http

import (
	"awesomeProject/metadata/internal/controller/metadata"
	"awesomeProject/metadata/internal/repository"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// Handler defines a movie metadata HTTP handler
type Handler struct {
	ctrl *metadata.Controller
}

// New creates a new movie metadata HTTP handler
func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// Getmetadata handles GET /metadata requests
func (h *Handler) GetMetadata(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	m, err := h.ctrl.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		log.Printf("repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("reponse encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
