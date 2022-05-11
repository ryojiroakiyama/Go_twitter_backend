package media

import (
	"io"
	"net/http"

	"yatter-backend-go/app/handler/httperror"
)

// Handle request for `POST /media`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()

	content, err := io.ReadAll(r.Body)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Write(content)
}

/* 以下のbodyの中のmultipart/form-dataを分割できるようなgoの機能を探す

------WebKitFormBoundaryWYmEcsXGLxLb50oE
Content-Disposition: form-data; name="file"; filename="Untitled.txt"
Content-Type: text/plain

akiyama content
------WebKitFormBoundaryWYmEcsXGLxLb50oE
Content-Disposition: form-data; name="description"

------WebKitFormBoundaryWYmEcsXGLxLb50oE--
*/
