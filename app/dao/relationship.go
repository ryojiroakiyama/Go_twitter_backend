package dao

import (
	"context"
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
func (r *relationship) Relationships(ctx context.Context, userID object.AccountID, targetID object.AccountID) (*object.Relationship, error) {
	schema := `
	SELECT`
	return nil, nil
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
