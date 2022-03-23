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
	status struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

// FindByAccountID : アカウントIDから投稿をとってくる
func (r *status) FindByAccountID(ctx context.Context, accountID object.AccountID) (*object.Status, error) {
	entity := new(object.Status)
	err := r.db.QueryRowxContext(ctx, "select * from status where account_id = ?", accountID).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}

	return entity, nil
}

// CreateStatus: アカウント作成
func (r *status) CreateStatus(ctx context.Context, entity *object.Status) error {
	schema := `insert into status (account_id, content) values (?, ?)`
	_, err := r.db.ExecContext(ctx, schema, entity.Account_ID, entity.Content)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
