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

	// parse multipart form
	account := auth.AccountOf(r)
	if dname := r.FormValue("display_name"); dname != "" {
		account.DisplayName = &dname
	}
	if note := r.FormValue("note"); note != "" {
		account.Note = &note
	}
	if filePath, err := formFileToFile(r, "avatar"); err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if filePath != "" {
		if account.Avatar != nil {
			os.Remove(*account.Avatar)
		}
		account.Avatar = &filePath
	}
	if filePath, err := formFileToFile(r, "header"); err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if filePath != "" {
		if account.Header != nil {
			os.Remove(*account.Header)
		}
		account.Header = &filePath
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

func formFileToFile(r *http.Request, key string) (string, error) {
	if file, _, err := r.FormFile(key); err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			return "", err
		}
	} else {
		path, err := fileio.WriteToTmpFile(file, "./.data/"+key, "")
		if err != nil {
			return "", err
		}
		return path, nil
	}
	return "", nil
}
