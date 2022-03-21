package accounts

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/handler/httperror"

	"github.com/go-chi/chi"
)

// Handle request for `GET /v1/accounts/{username}`
func (h *handler) Fetch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// handler 入力
	username := chi.URLParam(r, "username")

	// usercase 中(domainまで潜る)
	account, err := h.app.Dao.Account().FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
	}
	if account == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	// handler 出力
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
