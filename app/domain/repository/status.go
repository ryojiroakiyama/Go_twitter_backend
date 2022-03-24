package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Fetch status which has specified accountID
	FindByID(ctx context.Context, accountID object.AccountID) (*object.Status, error)

	// Create Status
	CreateStatus(ctx context.Context, entity *object.Status) error

	// Delete Status
	DeleteStatus(ctx context.Context, status_id object.StatusID, account_id object.AccountID) error
}
