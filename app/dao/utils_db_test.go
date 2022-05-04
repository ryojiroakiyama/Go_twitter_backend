package dao_test

//import (
//	"context"
//	"database/sql"
//	"errors"
//	"fmt"
//	"testing"
//	"yatter-backend-go/app/domain/object"

//	_ "github.com/go-sql-driver/mysql"
//	"github.com/jmoiron/sqlx"
//)

//type DBConfig interface {
//	FormatDSN() string
//}

//func NewDB(t *testing.T) *sqlx.DB {
//	t.Helper()
//	cfg := NewConfig(t)
//	db, err := sqlx.Open("mysql", cfg.FormatDSN())
//	if err != nil {
//		t.Fatal("Open: ", err)
//	}
//	pingErr := db.Ping()
//	if pingErr != nil {
//		t.Fatal("Ping: ", pingErr)
//	}

//	ctx := context.Background()
//	query := `
//	INSERT
//		INTO account
//		(username, password_hash) VALUES (?, ?)`
//	_, err = db.ExecContext(ctx, query, "rsudo", "rrr")
//	if err != nil {
//		t.Fatal("failt to insert", err)
//	}

//	account := new(object.Account)
//	query = `
//	SELECT *
//	FROM meta_account
//	WHERE username = ?`
//	err = db.QueryRowxContext(ctx, query, "rakiyama").StructScan(account)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			t.Fatal("no rows")
//		}
//		t.Fatal("fail to query", err)
//	}
//	fmt.Println("------->", account.Username)
//	return db
//}
