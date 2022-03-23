package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Fetch status which has specified accountID
	FindByAccountID(ctx context.Context, accountID object.AccountID) (*object.Status, error)

	// Create Status
	CreateStatus(ctx context.Context, entity *object.Status) error
}
