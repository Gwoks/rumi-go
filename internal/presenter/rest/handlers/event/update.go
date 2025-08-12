package event

import (
	"encoding/json"
	"net/http"
	"strconv"

	"rumi-go/internal/model"
)

// UpdateEvent handles HTTP requests for updating an event
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.URL.Query().Get("id")
	if eventIDStr == "" {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var req model.EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event, err := h.eventUsecase.UpdateEvent(r.Context(), eventID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}
