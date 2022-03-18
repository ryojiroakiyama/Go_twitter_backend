package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

// Handle request for `GET /v1/accounts/{username}`
func (h *handler) Fetch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := chi.URLParam(r, "username")

	accountRepository := h.app.Dao.Account()
	account, err := accountRepository.FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
	}
	if account == nil {
		httperror.BadRequest(w, fmt.Errorf("Account not found")) // 404(not found) is better?
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
