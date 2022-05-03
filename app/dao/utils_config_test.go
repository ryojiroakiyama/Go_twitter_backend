package dao_test

import (
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

func NewConfig() DBConfig {
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
