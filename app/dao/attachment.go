package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Attachment
	attachment struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewAttachment(db *sqlx.DB) repository.Attachment {
	return &attachment{db: db}
}

// FindByID : 指定IDのステータスの取得
func (r *attachment) FindByID(ctx context.Context, id object.AttachmentID) (*object.Attachment, error) {
	attachment := new(object.Attachment)
	query := `
	SELECT *
	FROM media
	WHERE id = ?`
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(attachment)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return attachment, nil
}

// Create: ステータス作成
func (r *attachment) Create(ctx context.Context, entity *object.Attachment) (object.AccountID, error) {
	query := `
	INSERT
		INTO attachment
		(type, url, description) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, entity.Type, entity.Url, entity.Description)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return id, nil
}

//// Delete: ステータス削除
//func (r *attachment) Delete(ctx context.Context, attachment_id object.AttachmentID, account_id object.AccountID) error {
//	query := `
//	DELETE
//	FROM attachment
//	WHERE id=? AND account_id=?`
//	_, err := r.db.ExecContext(ctx, query, attachment_id, account_id)
//	if err != nil {
//		return fmt.Errorf("%w", err)
//	}
//	return nil
//}
