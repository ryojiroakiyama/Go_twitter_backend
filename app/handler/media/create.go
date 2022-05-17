package media

import (
	"encoding/json"
	"net/http"
	"strings"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/fileio"
	"yatter-backend-go/app/handler/httperror"
)

// Handle request for `POST /media`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	media := new(object.Media)
	file, header, err := r.FormFile("file")
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	media.Url, err = fileio.WriteToTmpFile(file, "./.data/media", "")
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	media.Type = toMediaType(header.Header.Get("Content-Type"))
	if description := r.FormValue("description"); description != "" {
		media.Description = &description
	}

	id, err := h.app.Dao.Media().Create(ctx, media)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	media, err = h.app.Dao.Media().FindByID(ctx, id)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	} else if media == nil {
		httperror.LostAccount(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(media); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
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
