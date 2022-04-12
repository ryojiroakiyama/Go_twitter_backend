package accounts

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handle request for `POST /v1/accounts/{username}/unfollow`
func (h *handler) UnFollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	targetname, err := request.UserNameOf(r)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	target, err := h.app.Dao.Account().FindByUsername(ctx, targetname)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if target == nil {
		httperror.Error(w, http.StatusNotFound)
		return
	}

	user := auth.AccountOf(r)
	if user == nil {
		httperror.LostAccount(w)
		return
	}

	// delete follow
	err = h.app.Dao.Relationship().Delete(ctx, user.ID, target.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
	}

	relationship, err := h.app.Dao.Relationship().Relationship(ctx, user.ID, target.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(relationship); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
