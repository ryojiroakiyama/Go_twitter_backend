package statuses

import (
	"fmt"
	"net/http"

	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/httperror"
	"yatter-backend-go/app/handler/request"
)

// Handle request for `DELETE /v1/statuses/{id}`
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	account := auth.AccountOf(r)
	if account == nil {
		httperror.InternalServerError(w, fmt.Errorf("lost account"))
		return
	}

	id, err := request.IDOf(r)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	err = h.app.Dao.Status().DeleteStatus(ctx, id, account.ID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	//if err := json.NewEncoder(w).Encode(status); err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//}
}
