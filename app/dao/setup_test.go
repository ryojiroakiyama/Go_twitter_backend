package dao_test

import (
	"log"
	"os"
	"time"
	"yatter-backend-go/app/dao"

	"github.com/go-sql-driver/mysql"
)

func NewDao() dao.Dao {
	dao, err := dao.New(NewConfig())
	if err != nil {
		log.Fatal("NewTestDao() fail", err)
	}
	return dao
}

func NewConfig() dao.DBConfig {
	cfg := mysql.NewConfig()
	cfg.ParseTime = true
	if loc, err := time.LoadLocation(os.Getenv("TEST_MYSQL_TZ")); err != nil {
		log.Fatal("Invalid timezone")
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

func Done(d dao.Dao) {
	if err := d.InitAll(); err != nil {
		log.Fatal("InitAll fail")
	}
	if err := d.Close(); err != nil {
		log.Fatal("Close fail")
	}
}
