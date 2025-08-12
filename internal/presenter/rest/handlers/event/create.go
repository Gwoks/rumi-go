package event

import (
	"encoding/json"
	"net/http"

	"rumi-go/internal/model"
)

// CreateEvent handles HTTP requests for event creation
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req model.EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event, err := h.eventUsecase.CreateEvent(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}
