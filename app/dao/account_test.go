package dao_test

import (
	"context"
	"testing"
	"yatter-backend-go/app/domain/object"

	_ "github.com/go-sql-driver/mysql"
)

func TestAccount(t *testing.T) {
	dao := NewDao()
	a := object.Account{Username: "test"}
	dao.Account().Create(context.Background(), &a)
	defer Done(dao)
}
