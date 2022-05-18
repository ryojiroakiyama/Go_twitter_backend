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
	SELECT *
	FROM meta_account
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
	INSERT
		INTO account
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

// Create: アカウント更新
func (r *account) Update(ctx context.Context, account *object.Account) (error) {
	query := `
	UPDATE account
	SET    display_name = ?,
	       avatar = ?,
	       header = ?,
	       note = ?
	WHERE  id = ?`
	_, err := r.db.ExecContext(ctx, query, account.DisplayName, account.Avatar, account.Header, account.Note, account.ID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// Following: userがフォローしているアカウント集合を返す
func (r *account) Following(ctx context.Context, username string, limit int64) ([]object.Account, error) {
	var accounts []object.Account
	query := `
	SELECT
		ma.id,
		ma.username,
		ma.password_hash,
		ma.display_name,
		ma.avatar,
		ma.header,
		ma.note,
		ma.create_at,
		ma.following_count,
		ma.followers_count
	FROM
		meta_account AS ma
		INNER JOIN
		relationship AS r
		ON ma.id = r.follow_id
	WHERE r.user_id = (SELECT id FROM account WHERE username = ?)
	LIMIT ?`
	err := r.db.SelectContext(ctx, &accounts, query, username, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return accounts, nil
}

// Followers: userをフォロワーのアカウント集合を返す
func (r *account) Followers(ctx context.Context, username string, since_id int64, max_id int64, limit int64) ([]object.Account, error) {
	var accounts []object.Account
	query := `
	SELECT
		ma.id,
		ma.username,
		ma.password_hash,
		ma.display_name,
		ma.avatar,
		ma.header,
		ma.note,
		ma.create_at,
		ma.following_count,
		ma.followers_count
	FROM
		meta_account AS ma
		INNER JOIN
		relationship AS r
		ON ma.id = r.user_id
	WHERE
		r.follow_id = (SELECT id FROM account WHERE username = ?)
		AND ? <= ma.id
		AND ma.id <= ?
	LIMIT ?`
	err := r.db.SelectContext(ctx, &accounts, query, username, since_id, max_id, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return accounts, nil
}
