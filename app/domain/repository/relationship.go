package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Relationship interface {
	// Create Relationship to follow
	Create(ctx context.Context, userID object.AccountID, targetID object.AccountID) (object.RelationshipID, error)
}
