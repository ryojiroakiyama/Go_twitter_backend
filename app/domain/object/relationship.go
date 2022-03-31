package object

type (
	RelationshipID = int64

	Relationship struct {
		// The internal id of the rerationship
		ID RelationshipID `json:"-" db:"id"`

		// Target account id
		TargetID AccountID `json:"id,omitempty"`

		// Whether the user is currently following the account
		Following bool `json:"following,omitempty"`

		// Whether the user is currently being followed by the account
		FllowedBy bool `json:"followd_by,omitempty"`
	}
)
