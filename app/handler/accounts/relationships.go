package accounts

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/params"
)

// Handle request for `GET /v1/accounts/{username}/following`
func (h *handler) Relationships(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	targetNames := params.FromValueSplit(r, "username", ",")

	user := auth.AccountOf(r)
	if user == nil {
		httperror.LostAccount(w)
		return
	}

	var relationships []*object.Relationship
	for _, tname := range targetNames {
		target, err := h.app.Dao.Account().FindByUsername(ctx, tname)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}
		if target == nil {
			httperror.Error(w, http.StatusNotFound)
			return
		}
		r, err := h.app.Dao.Relationship().Fetch(ctx, user.ID, target.ID)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}
		relationships = append(relationships, r)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(relationships); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
