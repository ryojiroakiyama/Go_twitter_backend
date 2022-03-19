package statuses

import (
	"encoding/json"
	"fmt"
	"net/http"

	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
)

// Request body for `POST /v1/statuses`
type AddRequest struct {
	Status string
}

// tmp
type Status struct {
	Status string
}

//コンテキストからアカウントを取得して, その投稿を作成する
// Handle request for `POST /v1/statuses`
func (h *handler) Fetch(w http.ResponseWriter, r *http.Request) {
	/*accounts/create*/
	//ctx := r.Context()

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

	/*accounts/create*/
	//account := new(object.Account)
	//account.Username = req.Username
	//if err := account.SetPassword(req.Password); err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//}
	var s Status
	s.Status = req.Status

	/*accounts/create*/
	//if err := h.app.Dao.Account().CreateAccount(ctx, account); err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(s); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
