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
	status struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

//// FindByUsername : ユーザ名から
//func (r *status) FindByUsername(ctx context.Context, username string) (*object.Status, error) {
//	entity := new(object.Status)
//	err := r.db.QueryRowxContext(ctx, "select * from status where username = ?", username).StructScan(entity)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return nil, nil
//		}

//		return nil, fmt.Errorf("%w", err)
//	}

//	return entity, nil
//}

// CreateStatus: アカウント作成
func (r *status) CreateStatus(ctx context.Context, entity *object.Status) error {
	schema := `insert into status (account_id, content) values (?, ?)`
	_, err := r.db.ExecContext(ctx, schema, entity.Account_ID, entity.Content)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
