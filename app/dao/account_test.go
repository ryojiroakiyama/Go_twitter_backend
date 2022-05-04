package dao_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestAccount(t *testing.T) {
	NewDB(t)
}
