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

// IsFollowing: userがtargertをfollowしているかを取得する
func (r *relationship) IsFollowing(ctx context.Context, userID object.AccountID, targetID object.AccountID) (bool, error) {
	query := `
	SELECT id
	FROM relationship
	WHERE user_id =? AND follow_id=?`
	var id uint32
	if err := r.db.QueryRowContext(ctx, query, userID, targetID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Relationship: userとtargetのフォロー関係を取得する
func (r *relationship) Fetch(ctx context.Context, userID object.AccountID, targetID object.AccountID) (*object.Relationship, error) {
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

// Create: userがtargetをフォローする関係を登録
func (r *relationship) Create(ctx context.Context, userID object.AccountID, targetID object.AccountID) (object.RelationshipID, error) {
	query := `
	INSERT
		INTO relationship
		(user_id, follow_id) VALUES (?, ?)`
	result, err := r.db.ExecContext(ctx, query, userID, targetID)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return id, nil
}

// Delete: userがtargetをフォローする関係を削除
// followしてないアカウントを指定してもエラーにならない(UI部分ではそもそも表示されないこと想定)
func (r *relationship) Delete(ctx context.Context, userID object.AccountID, targetID object.AccountID) error {
	query := `
	DELETE
	FROM relationship
	WHERE user_id=? AND follow_id=?`
	_, err := r.db.ExecContext(ctx, query, userID, targetID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
