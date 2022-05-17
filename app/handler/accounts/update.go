package accounts

import (
	"errors"
	"fmt"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/fileio"
	"yatter-backend-go/app/handler/httperror"
)

// Handle request for `POST /v1/accounts/update_credentials`
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()

	// TODO: ファイルの切り出し, mediaのcreateと共有できるかも
	account := new(object.Account)
	if dname := r.FormValue("display_name"); dname != "" {
		account.DisplayName = &dname
	}
	if note := r.FormValue("note"); note != "" {
		account.Note = &note
	}
	if file, _, err := r.FormFile("avatar"); err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
		httperror.InternalServerError(w, err)
		return
		}
	} else {
		url, err := fileio.WriteToTmpFile(file, "./.data/avatar", "")
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}
		account.Avatar = &url
	}
	if file, _, err := r.FormFile("header"); err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
		httperror.InternalServerError(w, err)
		return
		}
	} else {
		url, err := fileio.WriteToTmpFile(file, "./.data/header", "")
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}
		account.Header = &url
	}
	fmt.Fprintf(w, "DisplayName: %v\n", *account.DisplayName)
	fmt.Fprintf(w, "Note: %v\n", *account.Note)
	fmt.Fprintf(w, "Avatar: %v\n", *account.Avatar)
	fmt.Fprintf(w, "Header: %v\n", *account.Header)

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
