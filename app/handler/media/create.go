package media

import (
	//"fmt"
	"fmt"
	"io"
	"os"

	//"mime"
	//"mime/multipart"
	"net/http"
	"strconv"

	//"strings"

	"yatter-backend-go/app/handler/httperror"
)

// TODO: リクエストをmedia objectにしてcreateする
// TODO: パースはr.FormFileで行けるかも
// Handle request for `POST /media`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()

	//body, err := io.ReadAll(r.Body)
	//if err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//}
	//w.Write(body)
	//w.Write([]byte("\n\n"))

	file, header, err := r.FormFile("file")
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	path, err := GenTmpFile(file, "./.data/media")
	w.Write([]byte(path))
	w.Write([]byte("\n"))
	w.Write([]byte("~~~header~~~\n"))
	w.Write([]byte("filename: " + header.Filename))
	w.Write([]byte("\n"))
	w.Write([]byte("size: " + strconv.FormatInt(header.Size, 10)))
	w.Write([]byte("\n"))
	w.Write([]byte("Content-Disposition: " + header.Header.Get("Content-Disposition")))
	w.Write([]byte("\n"))
	w.Write([]byte("Content-Type: " + header.Header.Get("Content-Type")))
	w.Write([]byte("\n"))

	//mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	//if err != nil {
	//	httperror.InternalServerError(w, err)
	//	return
	//}
	//w.Write([]byte(mediaType))
	//w.Write([]byte("\n"))
	//w.Write([]byte(params["boundary"]))
	//w.Write([]byte("\n"))
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
	//			return
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
func GenTmpFile(src io.Reader, dir string) (filePath string, err error) {
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
