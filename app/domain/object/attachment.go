package object

type (
	AttachmentID = int64

	Attachment struct {
		// The ID of the attachment
		ID AttachmentID `json:"id" db:"id"`

		// The content_type of the attachment
		Type string `json:"type" db:"type"`

		// Url(path) of the attachment
		Url string `json:"url,omitempty" db:"url"`

		// The time the attachment was created
		Description *string `json:"description,omitempty" db:"description"`
	}
)
