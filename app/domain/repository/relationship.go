package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Relationship interface {
	// Return whether the user is currently following the target
	IsFollowing(ctx context.Context, userID object.AccountID, targetID object.AccountID) (bool, error)

	// Fetch Relationship
	Relationship(ctx context.Context, userID object.AccountID, targetID object.AccountID) (*object.Relationship, error)

	// Create Relationship to follow
	Create(ctx context.Context, userID object.AccountID, targetID object.AccountID) (object.RelationshipID, error)

	// Fetch all following accounts
	FollowingAccounts(ctx context.Context, username string) ([]object.Account, error)

	// Fetch all follower accounts
	FollowerAccounts(ctx context.Context, username string) ([]object.Account, error)

	// Delete Relationship
	Delete(ctx context.Context, userID object.AccountID, targetID object.AccountID) error

	// Get number of following accounts
	NumberOfFollowingAccounts(ctx context.Context, username string) (int64, error)
}
