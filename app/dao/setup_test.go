package dao_test

import (
	"context"
	"os"
	"testing"
	"time"
	"yatter-backend-go/app/dao"

	"github.com/go-sql-driver/mysql"
	"yatter-backend-go/app/domain/object"
)

const (
	testUsername1 = "benben"
	testUsername2 = "jonjon"
	testUsername3 = "sonson"
)

func NewDao(t *testing.T) dao.Dao {
	t.Helper()
	dao, err := dao.New(NewConfig(t))
	if err != nil {
		t.Fatal("NewTestDao() fail", err)
	}
	return dao
}

func NewConfig(t *testing.T) dao.DBConfig {
	t.Helper()
	cfg := mysql.NewConfig()
	cfg.ParseTime = true
	if loc, err := time.LoadLocation(os.Getenv("TEST_MYSQL_TZ")); err != nil {
		t.Fatal("Invalid timezone")
	} else {
		cfg.Loc = loc
	}
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("TEST_MYSQL_HOST")
	cfg.User = os.Getenv("TEST_MYSQL_USER")
	cfg.Passwd = os.Getenv("TEST_MYSQL_PASSWORD")
	cfg.DBName = os.Getenv("TEST_MYSQL_DATABASE")
	return cfg
}

func CreateBaseTable(t *testing.T, dao dao.Dao) {
	t.Helper()
	ctx := context.Background()
	account1 := object.Account{Username: testUsername1}
	account2 := object.Account{Username: testUsername2}
	account3 := object.Account{Username: testUsername3}
	account1.ID = mustAccountCreate(t, dao, ctx, &account1)
	account2.ID = mustAccountCreate(t, dao, ctx, &account2)
	account3.ID = mustAccountCreate(t, dao, ctx, &account3)
	mustRelationshipCreate(t, dao, ctx, account1.ID, account2.ID)
	mustRelationshipCreate(t, dao, ctx, account1.ID, account3.ID)
	mustRelationshipCreate(t, dao, ctx, account2.ID, account3.ID)
}

func mustAccountCreate(t *testing.T, dao dao.Dao, ctx context.Context, ac *object.Account) int64 {
	t.Helper()
	id, err := dao.Account().Create(ctx, ac)
	if err != nil {
		t.Fatalf("BasaTable: %v", err)
	}
	return id
}

func mustRelationshipCreate(t *testing.T, dao dao.Dao, ctx context.Context, useID int64, targetID int64) int64 {
	t.Helper()
	id, err := dao.Relationship().Create(ctx, useID, targetID)
	if err != nil {
		t.Fatalf("BasaTable: %v", err)
	}
	return id
}
