package timelines

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/handler/httperror"
)

// Handle request for `GET /v1/timelines/public`
func (h *handler) FetchPublic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_ = r.FormValue("limit")

	statuses, err := h.app.Dao.TimeLine().GetAll(ctx)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if statuses == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
