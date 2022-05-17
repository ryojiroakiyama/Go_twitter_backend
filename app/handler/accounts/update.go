package accounts

import (
	"fmt"
	"net/http"
	"yatter-backend-go/app/handler/fileio"
	"yatter-backend-go/app/handler/httperror"
)

// Handle request for `POST /v1/accounts/update_credentials`
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	displayName := r.MultipartForm.Value["display_name"][0]
	note := r.MultipartForm.Value["note"][0]
	avatar, err := r.MultipartForm.File["avatar"][0].Open()
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	header, err := r.MultipartForm.File["header"][0].Open()
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	avatarURL, err := fileio.WriteToTmpFile(avatar, "./.data/media/avatar", "")
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	headerURL, err := fileio.WriteToTmpFile(header, "./.data/media/header", "")
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	fmt.Fprintf(w, "dname: %v, note: %v\n", displayName, note)
	fmt.Fprintf(w, "aurl: %v\n hurl: %v\n", avatarURL, headerURL)

	//var req requestSyntax
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	httperror.BadRequest(w, err)
	//	return
	//}

	//account := new(object.Account)
	//account.Username = req.Username
	//if err := account.SetPassword(req.Password); err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//}

	//if accountFound, err := h.app.Dao.Account().FindByUsername(ctx, account.Username); err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//} else if accountFound != nil {
	//	httperror.Error(w, http.StatusConflict)
	//	return
	//}

	//_, err := h.app.Dao.Account().Create(ctx, account)
	//if err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//}

	//account, err = h.app.Dao.Account().FindByUsername(ctx, account.Username)
	//if err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//} else if account == nil {
	//	httperror.LostAccount(w)
	//	return
	//}

	//w.Header().Set("Content-Type", "application/json")
	//if err := json.NewEncoder(w).Encode(account); err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//}
}
