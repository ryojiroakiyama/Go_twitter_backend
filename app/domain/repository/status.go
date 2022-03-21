package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Status interface {
	// Create Status
	CreateStatus(ctx context.Context, entity *object.Status) error
}
