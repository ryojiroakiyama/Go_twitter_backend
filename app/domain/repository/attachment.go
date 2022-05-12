package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Attachment interface {
	// Fetch attachment which has specified attachmentID
	FindByID(ctx context.Context, attachmentID object.Attachment) (*object.Attachment, error)

	//// Create Attachment
	//Create(ctx context.Context, entity *object.Attachment) (object.StatusID, error)

	//// Delete Attachment
	//Delete(ctx context.Context, status_id object.StatusID, account_id object.Attachment) error
}
