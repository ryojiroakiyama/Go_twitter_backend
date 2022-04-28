package accounts

import (
	"encoding/json"
	"math"
	"net/http"

	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/params"
	"yatter-backend-go/app/handler/request"
)

// Handle request for `GET /v1/accounts/{username}/followers`
func (h *handler) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	since_id := params.FormValue(r, "since_id", 0, 0, math.MaxInt64)
	max_id := params.FormValue(r, "max_id", math.MaxInt64, 0, math.MaxInt64)
	limit := params.FormValue(r, "limit", 40, 0, 80)

	username, err := request.UserNameOf(r)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	accounts, err := h.app.Dao.Account().Followers(ctx, username, since_id, max_id, limit)
	if err != nil {
		httperror.InternalServerError(w, err)
	}
	if accounts == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
