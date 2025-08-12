package event_report

import (
	"encoding/json"
	"net/http"
	"strconv"

	"rumi-go/internal/model"
)

// UpdateEventReport handles HTTP requests for updating an event report
func (h *EventReportHandler) UpdateEventReport(w http.ResponseWriter, r *http.Request) {
	reportIDStr := r.URL.Query().Get("id")
	if reportIDStr == "" {
		http.Error(w, "Report ID is required", http.StatusBadRequest)
		return
	}

	reportID, err := strconv.ParseInt(reportIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	var req model.EventReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	report, err := h.eventReportUsecase.UpdateEventReport(r.Context(), reportID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
