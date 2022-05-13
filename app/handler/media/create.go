package media

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"yatter-backend-go/app/handler/httperror"
)

/* 以下のbodyの中のmultipart/form-dataを分割できるようなgoの機能を探す

------WebKitFormBoundaryWYmEcsXGLxLb50oE
Content-Disposition: form-data; name="file"; filename="Untitled.txt"
Content-Type: text/plain

akiyama content
------WebKitFormBoundaryWYmEcsXGLxLb50oE
Content-Disposition: form-data; name="description"

------WebKitFormBoundaryWYmEcsXGLxLb50oE--
*/

// TODO: リクエストをattachment objectにしてcreateする
// TODO: パースはr.FormFileで行けるかも
// Handle request for `POST /media`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	w.Write([]byte(mediaType))
	w.Write([]byte("\n"))
	w.Write([]byte(params["boundary"]))
	w.Write([]byte("\n"))
	var content string
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(r.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				httperror.InternalServerError(w, err)
				return
			}
			slurp, err := io.ReadAll(p)
			if err != nil {
				httperror.InternalServerError(w, err)
				return
			}
			content = content + fmt.Sprintf("[Content-Disposition: %q] [Content-Type: %q] %q\n", p.Header.Get("Content-Disposition"), p.Header.Get("Content-Type"), slurp)
		}
	}
	w.Write([]byte(content))
}
