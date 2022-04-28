package accounts

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/params"
	"yatter-backend-go/app/handler/request"
)

// Handle request for `GET /v1/accounts/{username}/following`
func (h *handler) Following(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := params.FormValue(r, params.Limit, 40, 0, 80)

	username, err := request.UserNameOf(r)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	accounts, err := h.app.Dao.Account().Following(ctx, username, limit)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
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
