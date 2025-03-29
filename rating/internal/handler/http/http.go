package http

import (
	"awesomeProject/rating/internal/controller/rating"
	model "awesomeProject/rating/pkg"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	ctrl *rating.Controller
}

// Creates a new rating servive http handler
func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl}
}

// Handle handles PUT and GET /rating requests
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	recordID := model.RecordID(r.FormValue("id"))
	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recordType := model.RecordType(r.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		v, err := h.ctrl.GetAggregatedRating(r.Context(), recordID, recordType)
		if err != nil && errors.Is(err, rating.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("response encode error: %v\n", err)
		}
	case http.MethodPut:
		userID := model.UserID(r.FormValue("userId"))
		v, err := strconv.ParseFloat(r.FormValue("value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rating := &model.Rating{
			UserID: userID,
			Value:  model.RatingValue(v),
		}
		if err := h.ctrl.PutRating(r.Context(), recordID, recordType, rating); err != nil {
			log.Printf("repository put error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
