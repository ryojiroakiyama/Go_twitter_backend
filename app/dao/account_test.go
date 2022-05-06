package dao_test

import (
	"context"
	"testing"
	"yatter-backend-go/app/domain/object"

	_ "github.com/go-sql-driver/mysql"
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
	CreateBaseTable(t, dao)

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
				{
					Username: testUsername3,
				},
			},
			wantErr: false,
		},
		{
			name: "limit",
			args: args{
				ctx:      ctx,
				username: testUsername1,
				limit:    1,
			},
			want: []object.Account{
				{
					Username: testUsername2,
				},
			},
			wantErr: false,
		},
		{
			name: "no account",
			args: args{
				ctx:      ctx,
				username: testUsername3,
				limit:    10,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := dao.Account().Following(tt.args.ctx, tt.args.username, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("account.Following() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == len(tt.want) {
				for i, v := range got {
					if v.Username != tt.want[i].Username {
						t.Errorf("account.Following() = %v, want %v", got, tt.want)
					}
				}
			} else {
				t.Errorf("account.Following() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountFollowers(t *testing.T) {
	// set up
	dao := NewDao(t)
	defer dao.InitAll()
	ctx := context.Background()
	CreateBaseTable(t, dao)

	type args struct {
		ctx      context.Context
		username string
		since_id int64
		max_id   int64
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
				username: testUsername3,
				since_id: 0,
				max_id:   10,
				limit:    10,
			},
			want: []object.Account{
				{
					Username: testUsername1,
				},
				{
					Username: testUsername2,
				},
			},
			wantErr: false,
		},
		{
			name: "simple",
			args: args{
				ctx:      ctx,
				username: testUsername3,
				since_id: 0,
				max_id:   1,
				limit:    10,
			},
			want: []object.Account{
				{
					Username: testUsername1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := dao.Account().Followers(tt.args.ctx, tt.args.username, tt.args.since_id, tt.args.max_id, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("account.Followers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == len(tt.want) {
				for i, v := range got {
					if v.Username != tt.want[i].Username {
						t.Errorf("account.Followers() = %v, want %v", got, tt.want)
					}
				}
			} else {
				t.Errorf("account.Followers() = %v, want %v", got, tt.want)
			}
		})
	}
}
