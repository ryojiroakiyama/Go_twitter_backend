package media

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	media.Url, err = GenTmpFile("./.data/media", file)
	media.Type = header.Header.Get("Content-Type")
	w.Write([]byte("\n"))
	w.Write([]byte("type: " + media.Type))
	w.Write([]byte("\n"))
	w.Write([]byte("url: " + media.Url))
	w.Write([]byte("\n"))

	// ここでそんなのないって返ってくる, これを初めにしてもそうなる
	file, header, err = r.FormFile("description")
	if err != nil {
		w.Write([]byte("~~~~~~~~"))
		httperror.InternalServerError(w, err)
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if description := string(content); description != "" {
		media.Description = &description
	}
	w.Write([]byte("description: " + *media.Description))
	w.Write([]byte("\n"))

	//mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	//if err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//}
	//var content string
	//if strings.HasPrefix(mediaType, "multipart/") {
	//	mr := multipart.NewReader(r.Body, params["boundary"])
	//	for {
	//		p, err := mr.NextPart()
	//		if err == io.EOF {
	//			break
	//		}
	//		if err != nil {
	//			httperror.InternalServerError(w, err)
	//			returnf
	//		}
	//		slurp, err := io.ReadAll(p)
	//		if err != nil {
	//			httperror.InternalServerError(w, err)
	//			return
	//		}
	//		content = content + fmt.Sprintf("[Content-Disposition: %q] [Content-Type: %q] %q\n", p.Header.Get("Content-Disposition"), p.Header.Get("Content-Type"), slurp)
	//	}
	//}
	//w.Write([]byte(content))
}

//GenTmpFile creates a tmpprary file and write contents to the file.
//If successful, GenTmpFile returns a string
//which is the name of the created file and nil error.
//Else if faulse, GenTmpFile returns a empty string and any error encountered.
func GenTmpFile(dir string, src io.Reader) (filePath string, err error) {
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
