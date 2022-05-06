package dao_test

import (
	"context"
	"testing"
	"yatter-backend-go/app/domain/object"

	_ "github.com/go-sql-driver/mysql"
)

const (
	testUsername1 = "benben"
	testUsername2 = "sonson"
	testUsername3 = "jonjon"
)

func TestAccountFindByUsername(t *testing.T) {
	// set up
	dao := NewDao(t)
	defer dao.InitAll()
	ctx := context.Background()
	account := object.Account{Username: testUsername1}
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
				username: testUsername1,
			},
			want: &object.Account{
				Username: testUsername1,
			},
			wantErr: false,
		},
		{
			name: "no account",
			args: args{
				ctx:      ctx,
				username: "no such account",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := dao.Account().FindByUsername(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("account.FindByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if got.Username != tt.want.Username {
					t.Errorf("account.FindByUsername() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestAccountCreate(t *testing.T) {
	// set up
	dao := NewDao(t)
	defer dao.InitAll()
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		account *object.Account
	}
	tests := []struct {
		name    string
		args    args
		want    object.AccountID
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				ctx: ctx,
				account: &object.Account{
					Username: testUsername1,
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "account duplicate",
			args: args{
				ctx: ctx,
				account: &object.Account{
					Username: testUsername1,
				},
			},
			want:    1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := dao.Account().Create(tt.args.ctx, tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("account.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("account.Create() = %v, want %v", got, tt.want)
				return
			}
			ac, err := dao.Account().FindByUsername(ctx, tt.args.account.Username)
			if err != nil {
				t.Logf("find by username fail")
			} else if ac == nil {
				t.Errorf("the account was not created")
			}
		})
	}
}

func TestAccountFollowing(t *testing.T) {
	// set up
	dao := NewDao(t)
	defer dao.InitAll()
	ctx := context.Background()
	account1 := object.Account{Username: testUsername1}
	account2 := object.Account{Username: testUsername2}
	account3 := object.Account{Username: testUsername3}
	account1.ID, _ = dao.Account().Create(ctx, &account1)
	account2.ID, _ = dao.Account().Create(ctx, &account2)
	account3.ID, _ = dao.Account().Create(ctx, &account3)
	dao.Relationship().Create(ctx, account1.ID, account2.ID)

	type args struct {
		ctx      context.Context
		username string
		limit    int64
	}
	tests := []struct {
		name    string
		args    args
		want    []object.Account
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				ctx:      ctx,
				username: testUsername1,
				limit:    10,
			},
			want: []object.Account{
				{
					Username: testUsername2,
				},
			},
			wantErr: false,
		},
		//{
		//	name: "no account",
		//	args: args{
		//		ctx:      ctx,
		//		username: "no such account",
		//	},
		//	want:    nil,
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := dao.Account().Following(tt.args.ctx, tt.args.username, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("account.Following() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				for i, v := range got {
					if v.Username != tt.want[i].Username {
						t.Errorf("account.Following() = %v, want %v", got, tt.want)
					}
				}
			}
		})
	}
}
