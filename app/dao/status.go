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

// FindByID : アカウントIDから投稿をとってくる
func (r *status) FindByID(ctx context.Context, id object.StatusID) (*object.Status, error) {
	entity := new(object.Status)
	schema := `
	select 
		s.id, 
		s.content, 
		s.create_at, 
		a.id as "account.id", 
		a.username as "account.username",
		a.create_at as "account.create_at"
	from status as s inner join account as a 
	on s.account_id = a.id 
	where s.id = ?`
	err := r.db.QueryRowxContext(ctx, schema, id).StructScan(entity)
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
	_, err := r.db.ExecContext(ctx, schema, entity.Account.ID, entity.Content)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// DeleteStatus: アカウント削除
func (r *status) DeleteStatus(ctx context.Context, status_id object.StatusID, account_id object.AccountID) error {
	schema := `delete from status where id=? and account_id=?`
	_, err := r.db.ExecContext(ctx, schema, status_id, account_id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
