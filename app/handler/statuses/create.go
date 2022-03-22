package statuses

import (
	"encoding/json"
	"fmt"
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

// Response body for `POST /v1/statuses`
type AddResponse struct {
	AccountID        object.AccountID
	Account          object.Account
	Content          string
	CreateAt         object.DateTime
	MediaAttachments []AddMediaAttachments
}

// List of media attachments
type AddMediaAttachments struct {
	AccountID   object.AccountID
	Type        string
	Url         string
	Description string
}

// Handle request for `POST /v1/statuses`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	account := auth.AccountOf(r)
	if account == nil {
		httperror.InternalServerError(w, fmt.Errorf("lost account"))
		return
	}

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	status := new(object.Status)
	status.Content = req.Status
	status.Account_ID = account.ID

	if err := h.app.Dao.Status().CreateStatus(ctx, status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	res := new(AddResponse)
	res.AccountID = account.ID
	res.Account = *account
	res.Content = status.Content
	res.CreateAt = status.CreateAt
	//res.MediaAttachments: []AddMediaAttachments{
	//		{
	//			AccountID:   status.Account_ID,
	//			Type:        status.Type,
	//			Url:         *status.Url,
	//			Description: status.Description,
	//		},
	//	},
	//}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
