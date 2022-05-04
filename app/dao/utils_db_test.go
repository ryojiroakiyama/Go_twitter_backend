package dao_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"yatter-backend-go/app/domain/object"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DBConfig interface {
	FormatDSN() string
}

func NewDB(t *testing.T) {
	cfg := NewConfig(t)
	db, err := sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		t.Fatal("Open: ", err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		t.Fatal("Ping: ", pingErr)
	}

	ctx := context.Background()
	query := `
	INSERT
		INTO account
		(username, password_hash) VALUES (?, ?)`
	_, err = db.ExecContext(ctx, query, "rakiyama", "rrr")
	if err != nil {
		t.Fatal("failt to insert")
	}

	account := new(object.Account)
	query = `
	SELECT *
	FROM meta_account
	WHERE username = ?`
	err = db.QueryRowxContext(ctx, query, "rakiyama").StructScan(account)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			t.Fatal("no rows")
		}
		t.Fatal("fail to query")
	}
	fmt.Println("------->", account.Username)
}
