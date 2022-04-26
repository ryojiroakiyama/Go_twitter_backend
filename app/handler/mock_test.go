package handler_test

import (
	"context"
	"fmt"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"
)

type accountTableMock map[string]object.Account
type statusTableMock map[object.StatusID]object.Status
type relationshipTableMock map[object.RelationshipID]object.Relationship

// dbMock: 各テーブルをmapで模したdbモック
type dbMock struct {
	account      accountTableMock
	status       statusTableMock
	relationship relationshipTableMock
}

// accountMock: account repojitoryをモック
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

// statusMock: status repojitoryをモック
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

// relationshipMock: relationship repojitoryをモック
type relationshipMock struct {
	db *dbMock
}

func newRelationShipMock(db *dbMock) repository.Relationship {
	return &relationshipMock{db: db}
}

func (r *relationshipMock) IsFollowing(ctx context.Context, userID object.AccountID, targetID object.AccountID) (bool, error) {
	return true, nil
}

func (r *relationshipMock) Fetch(ctx context.Context, userID object.AccountID, targetID object.AccountID) (*object.Relationship, error) {
	isFollowing, err := r.IsFollowing(ctx, userID, targetID)
	if err != nil {
		return nil, err
	}
	isFollowed, err := r.IsFollowing(ctx, targetID, userID)
	if err != nil {
		return nil, err
	}
	return &object.Relationship{
		TargetID:  targetID,
		Following: isFollowing,
		FllowedBy: isFollowed,
	}, nil
}

func (r *relationshipMock) Create(ctx context.Context, userID object.AccountID, targetID object.AccountID) (object.RelationshipID, error) {
	return 0, nil
}

func (r *relationshipMock) Following(ctx context.Context, username string) ([]object.Account, error) {
	return nil, nil
}

func (r *relationshipMock) Followers(ctx context.Context, username string) ([]object.Account, error) {
	return nil, nil
}

func (r *relationshipMock) Delete(ctx context.Context, userID object.AccountID, targetID object.AccountID) error {
	return nil
}

// daoMock: daoのモック
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
	return newRelationShipMock(d.db)
}

func (d *daoMock) InitAll() error {
	d.db = nil
	return nil
}
