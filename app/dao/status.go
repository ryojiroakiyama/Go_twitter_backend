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

// Create: ステータス作成
func (r *status) Create(ctx context.Context, entity *object.Status) (object.AccountID, error) {
	schema := `insert into status (account_id, content) values (?, ?)`
	result, err := r.db.ExecContext(ctx, schema, entity.Account.ID, entity.Content)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return id, nil
}

// Delete: ステータス削除
func (r *status) Delete(ctx context.Context, status_id object.StatusID, account_id object.AccountID) error {
	schema := `delete from status where id=? and account_id=?`
	_, err := r.db.ExecContext(ctx, schema, status_id, account_id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// GetAll : ステータス情報を全て取得
func (r *status) All(ctx context.Context) ([]object.Status, error) {
	var entity []object.Status
	schema := `
	select
		s.id, 
		s.content, 
		s.create_at, 
		a.id as "account.id", 
		a.username as "account.username",
		a.create_at as "account.create_at"
	from status as s inner join account as a 
	on s.account_id = a.id`
	err := r.db.SelectContext(ctx, &entity, schema)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}

	return entity, nil
}
