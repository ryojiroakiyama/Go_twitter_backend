package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Account interface {
	// Fetch account which has specified username
	FindByUsername(ctx context.Context, username string) (*object.Account, error)

	// Create Account
	Create(ctx context.Context, entity *object.Account) (object.AccountID, error)

	// Update Account
	Update(ctx context.Context, entity *object.Account) (error)

	// Fetch all following accounts
	Following(ctx context.Context, username string, limit int64) ([]object.Account, error)

	// Fetch all follower accounts
	Followers(ctx context.Context, username string, since_id int64, max_id int64, limit int64) ([]object.Account, error)
}
