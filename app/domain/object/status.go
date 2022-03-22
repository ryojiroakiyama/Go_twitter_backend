package object

type (
	StatusID = int64

	// Status account
	Status struct {
		// The internal ID of the status
		ID StatusID `json:"-"`

		// The account ID connected with the status
		Account_ID AccountID `json:"id" db:"account_id"`

		// Contents of the status
		Content string `json:"content,omitempty" db:"content"`

		// The time the status was created
		CreateAt DateTime `json:"create_at,omitempty" db:"create_at"`

		// The type of content
		Type string `json:"type"`

		// URL to the content
		Url *string `json:"url"`

		// description of the status
		Description string `json:"description"`
	}
)
