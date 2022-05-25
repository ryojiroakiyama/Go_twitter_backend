package accounts

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

// Request body for `POST /v1/accounts`
type requestSyntax struct {
	Username string
	Password string
}

// Handle request for `POST /v1/accounts`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req requestSyntax
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	account := new(object.Account)
	account.Username = req.Username
	if err := account.SetPassword(req.Password); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if accountFound, err := h.app.Dao.Account().FindByUsername(ctx, account.Username); err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if accountFound != nil {
		httperror.Error(w, http.StatusConflict)
		return
	}

	_, err := h.app.Dao.Account().Create(ctx, account)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	account, err = h.app.Dao.Account().FindByUsername(ctx, account.Username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if account == nil {
		httperror.LostAccount(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
