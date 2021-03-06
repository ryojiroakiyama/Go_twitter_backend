package object

type (
	StatusID = int64

	Status struct {
		// The ID of the status
		ID StatusID `json:"id" db:"id"`

		// The Account of the status
		Account *Account `json:"account,omitempty" db:"account"`

		// Contents of the status
		Content string `json:"content,omitempty" db:"content"`

		// The time the status was created
		CreateAt DateTime `json:"create_at,omitempty" db:"create_at"`

		// Attachment's media id
		Media_ID *MediaID `json:"-" db:"media_id"`

		// Media attachment
		Attachment *Media `json:"media_attachments,omitempty" db:"media_attachments"`
	}
)
