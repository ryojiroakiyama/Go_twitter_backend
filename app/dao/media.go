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
	// Implementation for repository.Media
	media struct {
		db *sqlx.DB
	}
)

//NewMedia: Create media repository
func NewMedia(db *sqlx.DB) repository.Media {
	return &media{db: db}
}

//FindByID : 指定IDのmediaの取得
func (r *media) FindByID(ctx context.Context, id object.MediaID) (*object.Media, error) {
	media := new(object.Media)
	query := `
	SELECT *
	FROM media
	WHERE id = ?`
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(media)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	return media, nil
}

//Create: media作成
func (r *media) Create(ctx context.Context, entity *object.Media) (object.AccountID, error) {
	query := `
	INSERT
		INTO media
		(type, url, description) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, entity.Type, entity.Url, entity.Description)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return id, nil
}

//Delete: media削除
