package accounts

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	avatarURL, err := writeTmpFile("./.data/media/avatar", avatar)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	headerURL, err := writeTmpFile("./.data/media/header", header)
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

//writeTmpFile creates a tmpprary file and write contents to the file.
//If successful, writeTmpFile returns a file path and nil error, else returns error.
func writeTmpFile(dir string, src io.Reader) (filePath string, err error) {
	err = os.MkdirAll(dir, 0750)
	if err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("ToTmpFile: %v", err)
	}
	tmpfile, err := os.CreateTemp(dir, "")
	if err != nil {
		return "", fmt.Errorf("ToTmpFile: %v", err)
	}
	filePath = tmpfile.Name()
	defer func() {
		if cerr := tmpfile.Close(); cerr != nil {
			err = fmt.Errorf("ToTmpFile: %v", cerr)
		}
		if err != nil && filePath != "" {
			os.Remove(filePath)
		}
	}()
	_, err = io.Copy(tmpfile, src)
	if err != nil {
		return "", fmt.Errorf("ToTmpFile: %v", err)
	}
	err = tmpfile.Sync()
	return
}
