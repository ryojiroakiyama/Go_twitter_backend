package accounts

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"yatter-backend-go/app/handler/auth"
	"yatter-backend-go/app/handler/fileio"
	"yatter-backend-go/app/handler/httperror"
)

// Handle request for `POST /v1/accounts/update_credentials`
func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: 処理まとめる, 関数分ける
	account := auth.AccountOf(r)
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
		if account.Avatar != nil {
			os.Remove(*account.Avatar)
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
		if account.Header != nil {
			os.Remove(*account.Header)
		}
		account.Header = &url
	}

	if err := h.app.Dao.Account().Update(ctx, account); err != nil {
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
