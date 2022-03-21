package object

import (
//"fmt"

//"github.com/pkg/errors"
//"golang.org/x/crypto/bcrypt"
)

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

/* account */

//// Check if given password is match to account's password
//func (a *Account) CheckPassword(pass string) bool {
//	return bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(pass)) == nil
//}

//// Hash password and set it to account object
//func (a *Account) SetPassword(pass string) error {
//	passwordHash, err := generatePasswordHash(pass)
//	if err != nil {
//		return fmt.Errorf("generate error: %w", err)
//	}
//	a.PasswordHash = passwordHash
//	return nil
//}

//func generatePasswordHash(pass string) (PasswordHash, error) {
//	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
//	if err != nil {
//		return "", fmt.Errorf("hashing password failed: %w", errors.WithStack(err))
//	}
//	return PasswordHash(hash), nil
//}
