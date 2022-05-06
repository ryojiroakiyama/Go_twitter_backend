package dao_test

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestStatus_Fetch_Delete(t *testing.T) {
	// set up
	dao := NewDao(t)
	defer dao.InitAll()
	ctx := context.Background()
	CreateBaseTable(t, dao)

	if s, err := dao.Status().FindByID(ctx, 1); err != nil {
		t.Fatalf("Fetch %v", err)
	} else if s.Account.Username != testUsername1 {
		t.Errorf("Fetch returned: %v", s)
	}

	//if err := dao.Relationship().Delete(ctx, 1, 2); err != nil {
	//	t.Fatalf("Delete %v", err)
	//}

}
