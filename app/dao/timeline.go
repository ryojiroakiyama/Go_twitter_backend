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
	// Implementation for repository.Account
	timeline struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewTimeLine(db *sqlx.DB) repository.TimeLine {
	return &timeline{db: db}
}

// GetAll : ステータス情報を全て取得
func (r *timeline) GetAll(ctx context.Context) ([]object.Status, error) {
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
