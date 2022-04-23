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
	account struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewAccount(db *sqlx.DB) repository.Account {
	return &account{db: db}
}

// FindByUsername : ユーザ名からユーザを取得
func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	account := new(object.Account)
	query := `
	SELECT
		a.username,
		a.display_name,
		a.avatar,
		a.header,
		a.note,
		a.create_at,
		r_user.cnt_user_id_rows AS "following_count",
		r_follow.cnt_follow_id_rows AS "followers_count"
	FROM account AS a
		INNER JOIN
			(SELECT user_id, COUNT(*) AS cnt_user_id_rows FROM relationship GROUP BY user_id) AS r_user
			ON a.id = r_user.user_id
		INNER JOIN
			(SELECT follow_id, COUNT(*) AS cnt_follow_id_rows FROM relationship GROUP BY follow_id) AS r_follow
			ON a.id = r_follow.follow_id
	WHERE username = ?`
	err := r.db.QueryRowxContext(ctx, query, username).StructScan(account)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return account, nil
}

// Create: アカウント作成
func (r *account) Create(ctx context.Context, account *object.Account) (object.AccountID, error) {
	query := `
	INSERT INTO account
		(username, password_hash) VALUES (?, ?)`
	result, err := r.db.ExecContext(ctx, query, account.Username, account.PasswordHash)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return id, nil
}
