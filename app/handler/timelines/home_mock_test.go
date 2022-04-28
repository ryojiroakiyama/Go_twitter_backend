package timelines_test

import (
	"context"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"
)

const registeredUser = "john"

// 必要最低限の挙動を実装
type daoMock struct{}

func (d *daoMock) Account() repository.Account {
	return &accountMock{}
}
func (d *daoMock) Status() repository.Status {
	return &statusMock{}
}
func (d *daoMock) Relationship() repository.Relationship {
	return nil
}
func (d *daoMock) InitAll() error {
	return nil
}

type accountMock struct{}

func (r *accountMock) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	return &object.Account{
		Username: registeredUser,
	}, nil
}
func (r *accountMock) Create(ctx context.Context, account *object.Account) (object.AccountID, error) {
	return 0, nil
}
func (r *accountMock) Following(ctx context.Context, username string, limit int64) ([]object.Account, error) {
	return nil, nil
}
func (r *accountMock) Followers(ctx context.Context, username string, since_id int64, max_id int64, limit int64) ([]object.Account, error) {
	return nil, nil
}

type statusMock struct{}

func (r *statusMock) FindByID(ctx context.Context, id object.StatusID) (*object.Status, error) {
	return nil, nil
}
func (r *statusMock) Create(ctx context.Context, Status *object.Status) (object.AccountID, error) {
	return 0, nil
}
func (r *statusMock) Delete(ctx context.Context, status_id object.StatusID, account_id object.AccountID) error {
	return nil
}
func (r *statusMock) AllStatuses(ctx context.Context, since_id int64, max_id int64, limit int64) ([]object.Status, error) {
	return nil, nil
}
func (r *statusMock) FollowingStatuses(ctx context.Context, username string, since_id int64, max_id int64, limit int64) ([]object.Status, error) {
	return []object.Status{}, nil
}
