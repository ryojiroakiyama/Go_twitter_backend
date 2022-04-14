package timelines

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

// Handle request for `GET /v1/timelines/home`
func (h *handler) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_ = r.FormValue("only_media")
	_ = r.FormValue("max_id")
	_ = r.FormValue("since_id")
	_ = r.FormValue("limit")

	account := auth.AccountOf(r)
	if account == nil {
		httperror.LostAccount(w)
		return
	}
	statuses, err := h.app.Dao.Status().FollowingStatuses(ctx, account.Username)
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
