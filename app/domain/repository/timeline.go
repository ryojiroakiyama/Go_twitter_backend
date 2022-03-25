package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type TimeLine interface {
	// Fetch all statuses
	GetAll(ctx context.Context) ([]object.Status, error)
}
