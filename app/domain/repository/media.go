package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Media interface {
	// Fetch media which has specified mediaID
	FindByID(ctx context.Context, id object.MediaID) (*object.Media, error)

	// Create Media
	Create(ctx context.Context, entity *object.Media) (object.AccountID, error)

	//// Delete Media
	//Delete(ctx context.Context, status_id object.StatusID, account_id object.Media) error
}
