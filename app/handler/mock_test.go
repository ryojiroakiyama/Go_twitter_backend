package handler_test

import (
	"context"
	"fmt"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"
)

type accountTableMock map[string]object.Account
type statusTableMock map[object.StatusID]object.Status

type dbMock struct {
	account accountTableMock
	status  statusTableMock
}

// accountMock: accountに関するrepojitoryとtableをモック
type accountMock struct {
	db *dbMock
}

func newAccountMock(db *dbMock) repository.Account {
	return &accountMock{db: db}
}

func (r *accountMock) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	a, exist := r.db.account[username]
	if !exist {
		return nil, nil
	}
	return &a, nil
}

func (r *accountMock) Create(ctx context.Context, entity *object.Account) (object.AccountID, error) {
	_, exist := r.db.account[entity.Username]
	if exist {
		return 0, fmt.Errorf("Create: Account aready exist")
	}
	r.db.account[entity.Username] = *entity
	id := len(r.db.account) + 1
	return int64(id), nil
}

// statusMock: statusに関するrepojitoryとtableをモック
type statusMock struct {
	db *dbMock
}

func newStatusMock(db *dbMock) repository.Status {
	return &statusMock{db: db}
}

func (r *statusMock) FindByID(ctx context.Context, id object.StatusID) (*object.Status, error) {
	s, exist := r.db.status[id]
	if !exist {
		return nil, nil
	}
	return &s, nil
}

func (r *statusMock) Create(ctx context.Context, entity *object.Status) (object.AccountID, error) {
	entity.ID = int64(len(r.db.status) + 1)
	r.db.status[entity.ID] = *entity
	return entity.ID, nil
}

func (r *statusMock) Delete(ctx context.Context, status_id object.StatusID, account_id object.AccountID) error {
	s, exist := r.db.status[status_id]
	if !exist {
		return fmt.Errorf("Delete: No status matched")
	} else if s.Account.ID != account_id {
		return fmt.Errorf("Delete: No status matched")
	}
	delete(r.db.status, status_id)
	return nil
}

func (r *statusMock) AllStatuses(ctx context.Context) ([]object.Status, error) {
	var statuses []object.Status
	for _, value := range r.db.status {
		statuses = append(statuses, value)
	}
	return statuses, nil
}

func (r *statusMock) FollowingStatuses(ctx context.Context, username string) ([]object.Status, error) {
	return nil, nil
}

type daoMock struct {
	db *dbMock
}

func (d *daoMock) Account() repository.Account {
	return newAccountMock(d.db)
}

func (d *daoMock) Status() repository.Status {
	return newStatusMock(d.db)
}

func (d *daoMock) Relationship() repository.Relationship {
	return nil
}

func (d *daoMock) InitAll() error {
	d.db = nil
	return nil
}
