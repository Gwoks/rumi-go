package event_report

import (
	"net/http"
	"strconv"
)

// DeleteEventReport handles HTTP requests for deleting an event report
func (h *EventReportHandler) DeleteEventReport(w http.ResponseWriter, r *http.Request) {
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

	if err := h.eventReportUsecase.DeleteEventReport(r.Context(), reportID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
