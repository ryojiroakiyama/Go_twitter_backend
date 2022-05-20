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

// FindByID : 指定IDのステータスの取得
func (r *status) FindByID(ctx context.Context, id object.StatusID) (*object.Status, error) {
	status := new(object.Status)
	query := `
	SELECT
		s.id,
		s.content,
		s.create_at,
		s.media_id,
		ma.id AS "account.id",
		ma.username AS "account.username",
		ma.password_hash AS "account.password_hash",
		ma.display_name AS "account.display_name",
		ma.avatar AS "account.avatar",
		ma.header AS "account.header",
		ma.note AS "account.note",
		ma.create_at AS "account.create_at",
		ma.following_count AS "account.following_count",
		ma.followers_count AS "account.followers_count"
	FROM
		status AS s
		INNER JOIN
		meta_account AS ma
		ON s.account_id = ma.id
	WHERE s.id = ?`
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	if err = r.insertMedia(ctx, status); err != nil {
		return nil, err
	}
	return status, nil
}

// Create: ステータス作成
func (r *status) Create(ctx context.Context, entity *object.Status) (object.StatusID, error) {
	if entity.Attachment != nil {
		return r.create(ctx, entity)
	} else {
		return r.createNoAttachment(ctx, entity)
	}
}

// Create: attachmentありのステータス作成
func (r *status) create(ctx context.Context, entity *object.Status) (object.StatusID, error) {
	query := `
		INSERT
			INTO status
			(account_id, content, media_id) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, entity.Account.ID, entity.Content, entity.Attachment.ID)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return id, nil
}

// Create: attachmentなしのステータス作成
func (r *status) createNoAttachment(ctx context.Context, entity *object.Status) (object.StatusID, error) {
	query := `
	INSERT
		INTO status
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

// AllStatuses : ステータス情報を全て取得
func (r *status) AllStatuses(ctx context.Context, since_id int64, max_id int64, limit int64) ([]object.Status, error) {
	var statuses []object.Status
	query := `
	SELECT
		s.id,
		s.content,
		s.create_at,
		s.media_id,
		ma.id AS "account.id",
		ma.username AS "account.username",
		ma.password_hash AS "account.password_hash",
		ma.display_name AS "account.display_name",
		ma.avatar AS "account.avatar",
		ma.header AS "account.header",
		ma.note AS "account.note",
		ma.create_at AS "account.create_at",
		ma.following_count AS "account.following_count",
		ma.followers_count AS "account.followers_count"
	FROM
		status AS s
		INNER JOIN
		meta_account AS ma 
		ON s.account_id = ma.id
	WHERE
			? <= s.id
		AND s.id <= ?
	LIMIT ?`
	err := r.db.SelectContext(ctx, &statuses, query, since_id, max_id, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	for i := 0; i < len(statuses); i++ {
		r.insertMedia(ctx, &statuses[i])
	}
	return statuses, nil
}

func (r *status) RelationStatuses(ctx context.Context, user_id object.AccountID, since_id int64, max_id int64, limit int64) ([]object.Status, error) {
	var statuses []object.Status
	// メインクエリ	: サブクエリテーブルとstatusテーブルをJOIN
	// サブクエリ	: userがフォローしているアカウントとuser自身のアカウントのみで構成されたmeta_accountテーブル
	//              現状: 複数からフォローされている場合にダブるのでGROUP BYしてる..
	//              前回: INNER JOIN していたらuserが誰にもフォローされていない場合拾えない
	query := `
	SELECT
		s.id, 
		s.content, 
		s.create_at,
		s.media_id,
		ma.id AS "account.id",
		ma.username AS "account.username",
		ma.password_hash AS "account.password_hash",
		ma.display_name AS "account.display_name",
		ma.avatar AS "account.avatar",
		ma.header AS "account.header",
		ma.note AS "account.note",
		ma.create_at AS "account.create_at",
		ma.following_count AS "account.following_count",
		ma.followers_count AS "account.followers_count"
	FROM
		status AS s
		INNER JOIN
			(SELECT
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
				LEFT OUTER JOIN
				relationship AS r
				ON ma.id = r.follow_id
			WHERE
				r.user_id = ?
				OR ma.id = ?
			GROUP BY
				ma.id)
		AS ma
		ON s.account_id = ma.id
	WHERE
			? <= s.id
		AND s.id <= ?
	LIMIT ?`
	err := r.db.SelectContext(ctx, &statuses, query, user_id, user_id, since_id, max_id, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w", err)
	}
	for i := 0; i < len(statuses); i++ {
		r.insertMedia(ctx, &statuses[i])
	}
	return statuses, nil
}

//insertMedia inserts media into status.Attachment if status.Media_ID != nil
// クエリ一発でattachmentありorなし含めてstatus情報を取ってくるのが思いつかなかったので,
// status情報を取ってきた後にこの関数でmediaを取ってくる
func (r *status) insertMedia(ctx context.Context, statuse *object.Status) error {
	if statuse.Media_ID != nil {
		media := new(object.Media)
		query := `
			SELECT *
			FROM media
			WHERE id = ?`
		err := r.db.QueryRowxContext(ctx, query, statuse.Media_ID).StructScan(media)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("status has media that doesn't exist")
			}
			return fmt.Errorf("%w", err)
		}
		statuse.Attachment = media
	}
	return nil
}
