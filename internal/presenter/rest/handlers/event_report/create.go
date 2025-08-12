package event_report

import (
	"encoding/json"
	"net/http"

	"rumi-go/internal/model"
)

// CreateEventReport handles HTTP requests for event report creation
func (h *EventReportHandler) CreateEventReport(w http.ResponseWriter, r *http.Request) {
	var req model.EventReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	report, err := h.eventReportUsecase.CreateEventReport(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(report)
}
