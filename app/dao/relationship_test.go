package dao_test

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestRelationship(t *testing.T) {
	// set up
	dao := NewDao(t)
	defer dao.InitAll()
	ctx := context.Background()
	CreateBaseTable(t, dao)

	if r, err := dao.Relationship().Fetch(ctx, 1, 2); err != nil {
		t.Fatalf("Fetch %v", err)
	} else if !(r.Following && !r.FllowedBy) {
		t.Errorf("Fetch returned: %v", r)
	}

	if err := dao.Relationship().Delete(ctx, 1, 2); err != nil {
		t.Fatalf("Delete %v", err)
	}

	if r, err := dao.Relationship().Fetch(ctx, 1, 2); err != nil {
		t.Fatalf("Fetch %v", err)
	} else if !(!r.Following && !r.FllowedBy) {
		t.Errorf("Fetch returned: %v", r)
	}
}
