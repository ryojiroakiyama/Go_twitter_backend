package dao_test

import (
	"context"
	"testing"
	"yatter-backend-go/app/domain/object"

	_ "github.com/go-sql-driver/mysql"
)

const (
	testuser = "benben"
)

func TestAccountFindByUsername(t *testing.T) {
	// set up
	dao := NewDao(t)
	defer dao.InitAll()
	ctx := context.Background()
	account := object.Account{Username: testuser}
	dao.Account().Create(ctx, &account)

	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    *object.Account
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				ctx:      ctx,
				username: testuser,
			},
			want: &object.Account{
				Username: testuser,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := dao.Account().FindByUsername(ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("account.FindByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Username != tt.want.Username {
				t.Errorf("account.FindByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
