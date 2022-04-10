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
	query := `
	SELECT
		s.id,
		s.content,
		s.create_at,
		a.id AS "account.id",
		a.username AS "account.username",
		a.create_at AS "account.create_at"
	FROM status AS s INNER JOIN account AS a
		ON s.account_id = a.id
	WHERE s.id = ?`
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}

// Create: ステータス作成
func (r *status) Create(ctx context.Context, entity *object.Status) (object.StatusID, error) {
	query := `
	INSERT INTO status
		(account_id, content) VALUES (?, ?)`
	result, err := r.db.ExecContext(ctx, query, entity.Account.ID, entity.Content)
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
	query := `
	DELETE
	FROM status
	WHERE id=? AND account_id=?`
	_, err := r.db.ExecContext(ctx, query, status_id, account_id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// GetAll : ステータス情報を全て取得
func (r *status) All(ctx context.Context) ([]object.Status, error) {
	var entity []object.Status
	query := `
	SELECT
		s.id, 
		s.content, 
		s.create_at, 
		a.id AS "account.id", 
		a.username AS "account.username",
		a.create_at AS "account.create_at"
	FROM status AS s INNER JOIN account AS a 
		ON s.account_id = a.id`
	err := r.db.SelectContext(ctx, &entity, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return entity, nil
}
