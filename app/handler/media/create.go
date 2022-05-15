package media

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/httperror"
)

// Handle request for `POST /media`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()

	media := new(object.Media)
	file, header, err := r.FormFile("file")
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	media.Url, err = writeTmpFile("./.data/media", file)
	media.Type = toMediaType(header.Header.Get("Content-Type"))
	w.Write([]byte("\n"))
	w.Write([]byte("type: " + media.Type))
	w.Write([]byte("\n"))
	w.Write([]byte("url: " + media.Url))
	w.Write([]byte("\n"))

	if description := r.FormValue("description"); description != "" {
		media.Description = &description
	}
	w.Write([]byte("description: " + *media.Description))
	w.Write([]byte("\n"))
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

// Media objectが対応するtypeに変換
func toMediaType(fileType string) string {
	var mtype string
	for _, mtype = range object.MediaType {
		if strings.HasPrefix(fileType, mtype) {
			return mtype
		}
	}
	return mtype
}
