package http

import (
	"awesomeProject/movie/internal/controller/movie"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Handler struct {
	ctrl *movie.Controller
}

// Creates a new movie handler
func New(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl}
}

// Get movie details handles Get /movie rquests
func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	details, err := h.ctrl.Get(req.Context(), id)
	if err != nil && errors.Is(err, movie.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("response encode error: %v\n", err)
	}
}
