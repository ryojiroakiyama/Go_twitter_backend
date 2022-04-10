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
	query := `
	INSERT INTO relationship
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

func (r *relationship) FollowingAccounts(ctx context.Context, username string) ([]object.Account, error) {
	var accounts []object.Account
	query := `
	SELECT
		a.id,
		a.username,
		a.password_hash,
		a.display_name,
		a.avatar,
		a.header,
		a.note,
		a.create_at
	FROM account AS a INNER JOIN relationship AS r
		ON a.id = r.follow_id
	WHERE r.user_id = (SELECT id FROM account WHERE username = ?)`
	err := r.db.SelectContext(ctx, &accounts, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return accounts, nil
}

func (r *relationship) FollowerAccounts(ctx context.Context, username string) ([]object.Account, error) {
	var accounts []object.Account
	query := `
	SELECT
		a.id,
		a.username,
		a.password_hash,
		a.display_name,
		a.avatar,
		a.header,
		a.note,
		a.create_at
	FROM account AS a INNER JOIN relationship AS r
		ON a.id = r.user_id
	WHERE r.follow_id = (SELECT id FROM account WHERE username = ?)`
	err := r.db.SelectContext(ctx, &accounts, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return accounts, nil
}
