package dao

import (
	"context"
	"database/sql"
	"fmt"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Status
	relationship struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewRelationship(db *sqlx.DB) repository.Relationship {
	return &relationship{db: db}
}

// Fetch: Relationship オブジェクト返す
func (r *relationship) IsFollowing(ctx context.Context, userID object.AccountID, targetID object.AccountID) (bool, error) {
	schema := `
	SELECT *
	FROM relationship
	WHERE user_id =? AND follow_id=?`
	var id uint32
	var uid uint32
	var fid uint32
	if err := r.db.QueryRowContext(ctx, schema, userID, targetID).Scan(&id, &uid, &fid); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Create: フォロー関係作成
func (r *relationship) Create(ctx context.Context, userID object.AccountID, targetID object.AccountID) (object.RelationshipID, error) {
	schema := `
	INSERT INTO relationship
		(user_id, follow_id) VALUES (?, ?)`
	result, err := r.db.ExecContext(ctx, schema, userID, targetID)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return id, nil
}
