package dao_test

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestStatus_FindByID_Delete(t *testing.T) {
	// set up
	dao := NewDao(t)
	defer dao.InitAll()
	ctx := context.Background()
	CreateBaseTable(t, dao)

	if s, err := dao.Status().FindByID(ctx, 1); err != nil {
		t.Fatalf("FindByID %v", err)
	} else if s.Account.Username != testUsername1 {
		t.Errorf("FindByID returned: %v", s)
	}

	if err := dao.Status().Delete(ctx, 1, 1); err != nil {
		t.Fatalf("Delete %v", err)
	}

	if s, err := dao.Status().FindByID(ctx, 1); err != nil {
		t.Fatalf("FindByID %v", err)
	} else if s != nil {
		t.Errorf("FindByID returned: %v", s)
	}
}
