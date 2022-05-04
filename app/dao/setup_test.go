package dao_test

import (
	"os"
	"testing"
	"time"
	"yatter-backend-go/app/dao"

	"github.com/go-sql-driver/mysql"
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
