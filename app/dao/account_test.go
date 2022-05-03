package dao_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"
	"yatter-backend-go/app/domain/object"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DBConfig interface {
	FormatDSN() string
}

func TestAccount(t *testing.T) {
	cfg := NewConfig()
	db, err := sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("Open: ", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal("Ping: ", pingErr)
	}

	account := new(object.Account)
	ctx := context.Background()
	query := `
	SELECT *
	FROM meta_account
	WHERE username = ?`
	err = db.QueryRowxContext(ctx, query, "rakiyama").StructScan(account)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Fatal("no rows")
		}
		log.Fatal("fail to query")
	}
	fmt.Println("------->", account.Username)
}
