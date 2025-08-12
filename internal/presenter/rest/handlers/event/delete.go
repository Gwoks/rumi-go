package event

import (
	"net/http"
	"strconv"
)

// DeleteEvent handles HTTP requests for deleting an event
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
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

	if err := h.eventUsecase.DeleteEvent(r.Context(), eventID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
