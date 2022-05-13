package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Attachment interface {
	// Fetch attachment which has specified attachmentID
	FindByID(ctx context.Context, id object.AttachmentID) (*object.Attachment, error)

	// Create Attachment
	Create(ctx context.Context, entity *object.Attachment) (object.AccountID, error)

	//// Delete Attachment
	//Delete(ctx context.Context, status_id object.StatusID, account_id object.Attachment) error
}
