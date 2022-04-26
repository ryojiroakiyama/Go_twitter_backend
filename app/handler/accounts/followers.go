package accounts

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handle request for `GET /v1/accounts/{username}/followers`
func (h *handler) Followers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_ = r.FormValue("max_id")
	_ = r.FormValue("since_id")
	_ = r.FormValue("limit")

	username, err := request.UserNameOf(r)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	accounts, err := h.app.Dao.Account().Followers(ctx, username)
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
