package statuses

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

// Request body for `POST /v1/statuses`
type AddRequest struct {
	Status    string
	Media_ids []int64
}

// Handle request for `POST /v1/statuses`
// TODO: multiple attachments
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	account := auth.AccountOf(r)
	if account == nil {
		httperror.LostAccount(w)
		return
	}

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	status := new(object.Status)
	status.Content = req.Status
	status.Account = account
	if len(req.Media_ids) != 0 {
		attachment, err := h.app.Dao.Media().FindByID(ctx, req.Media_ids[0])
		if err != nil {
			httperror.LostAccount(w)
			return
		}
		status.Attachment = attachment
	}

	id, err := h.app.Dao.Status().Create(ctx, status)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	status, err = h.app.Dao.Status().FindByID(ctx, id)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
