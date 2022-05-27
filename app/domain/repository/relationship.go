package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

//Relationship inplements db operation abount follow relationships
type Relationship interface {
	// Return whether the user is currently following the target
	IsFollowing(ctx context.Context, userID object.AccountID, targetID object.AccountID) (bool, error)

	// Fetch infomation of relationship
	Fetch(ctx context.Context, userID object.AccountID, targetID object.AccountID) (*object.Relationship, error)

	// Create relationship
	Create(ctx context.Context, userID object.AccountID, targetID object.AccountID) (object.RelationshipID, error)

	// Delete relationship
	Delete(ctx context.Context, userID object.AccountID, targetID object.AccountID) error
}
