package timelines

import (
	//"encoding/json"
	"net/http"
	//"yatter-backend-go/app/handler/httperror"
	//"yatter-backend-go/app/handler/request"
)

// Handle request for `GET /v1/timelines/home`
func (h *handler) FetchHome(_ http.ResponseWriter, _ *http.Request) {
	//ctx := r.Context()

	//id, err := request.IDOf(r)
	//if err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//}

	//status, err := h.app.Dao.Status().FindByID(ctx, id)
	//if err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//}
	//if status == nil {
	//	httperror.Error(w, http.StatusNotFound)
	//	return
	//}

	//w.Header().Set("Content-Type", "application/json")
	//if err := json.NewEncoder(w).Encode(status); err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//}
}
