package timelines

import (
	"encoding/json"
	"math"
	"net/http"

	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/params"
)

// Handle request for `GET /v1/timelines/public`
func (h *handler) Public(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_ = params.FormValue(r, params.OnlyMedia, 40, 0, 80)
	since_id := params.FormValue(r, params.SinceID, 0, 0, math.MaxInt64)
	max_id := params.FormValue(r, params.MaxID, math.MaxInt64, 0, math.MaxInt64)
	limit := params.FormValue(r, params.Limit, 40, 0, 80)

	statuses, err := h.app.Dao.Status().AllStatuses(ctx, since_id, max_id, limit)
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
