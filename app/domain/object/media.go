package object

var MediaType = [...]string{"image", "video", "gifv", "unknown"}

type (
	MediaID = int64

	Media struct {
		// The ID of the media
		ID MediaID `json:"id" db:"id"`

		// The content_type of the media, ["image", "video", "gifv", "unknown"]
		Type string `json:"type" db:"type"`

		// Url(path) of the media
		Url string `json:"url,omitempty" db:"url"`

		// The time the media was created
		Description *string `json:"description,omitempty" db:"description"`
	}
)
