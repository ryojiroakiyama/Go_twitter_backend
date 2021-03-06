package accounts

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handle request for `POST /v1/accounts/{username}/follow`
func (h *handler) Follow(w http.ResponseWriter, r *http.Request) {
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

	// create follow
	if alreadyFollowing, err := h.app.Dao.Relationship().IsFollowing(ctx, user.ID, target.ID); err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if !alreadyFollowing && user.ID != target.ID {
		_, err = h.app.Dao.Relationship().Create(ctx, user.ID, target.ID)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}
	}

	relationship, err := h.app.Dao.Relationship().Fetch(ctx, user.ID, target.ID)
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
