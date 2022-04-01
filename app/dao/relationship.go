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

// IsFollowing: followしているかどうかを返す
func (r *relationship) IsFollowing(ctx context.Context, userID object.AccountID, targetID object.AccountID) (bool, error) {
	schema := `
	SELECT *
	FROM relationship
	WHERE user_id =? AND follow_id=?`
	var id, uid, fid uint32
	if err := r.db.QueryRowContext(ctx, schema, userID, targetID).Scan(&id, &uid, &fid); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Relationship: フォロー関係を返す
func (r *relationship) Relationship(ctx context.Context, userID object.AccountID, targetID object.AccountID) (*object.Relationship, error) {
	isFollowing, err := r.IsFollowing(ctx, userID, targetID)
	if err != nil {
		return nil, err
	}
	isFollowed, err := r.IsFollowing(ctx, targetID, userID)
	if err != nil {
		return nil, err
	}
	return &object.Relationship{
		TargetID:  targetID,
		Following: isFollowing,
		FllowedBy: isFollowed,
	}, nil
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
