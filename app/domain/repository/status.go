package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Fetch status which has specified accountID
	FindByID(ctx context.Context, accountID object.AccountID) (*object.Status, error)

	// Create Status
	Create(ctx context.Context, entity *object.Status) (object.StatusID, error)

	// Delete Status
	Delete(ctx context.Context, status_id object.StatusID, account_id object.AccountID) error

	// Fetch all statuses
	AllStatuses(ctx context.Context) ([]object.Status, error)

	// Fetch following account's statuses
	FollowingStatuses(ctx context.Context, username string) ([]object.Status, error)
}
