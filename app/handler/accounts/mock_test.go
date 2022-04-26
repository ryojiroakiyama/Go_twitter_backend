package accounts_test

import (
	"context"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"
)

/*
 * dao, db, repojitoryをモック
 * 構成はdaoと同じようにした
 */

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

// dbMock: 各テーブルをmapで模したdbモック
type dbMock struct {
	account      accountTableMock
	status       statusTableMock
	relationship relationshipTableMock
}

type accountTableMock map[int64]accountData
type accountData struct {
	id       int64
	username string
}

type statusTableMock map[int64]statusData
type statusData struct {
	id       int64
	username string
}

type relationshipTableMock map[int64]relationshipData
type relationshipData struct {
	userID   int64
	targetID int64
}

// accountMock: account repojitoryをモック
type accountMock struct {
	db *dbMock
}

func newAccountMock(db *dbMock) repository.Account {
	return &accountMock{db: db}
}

func (r *accountMock) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	for _, value := range r.db.account {
		if value.username == username {
			return &object.Account{
					ID:       value.id,
					Username: value.username},
				nil
		}
	}
	return nil, nil
}

func (r *accountMock) Create(ctx context.Context, account *object.Account) (object.AccountID, error) {
	newID := int64(len(r.db.account) + 1)
	r.db.account[int64(newID)] = accountData{
		id:       newID,
		username: account.Username}
	return newID, nil
}

func (r *accountMock) Following(ctx context.Context, username string) ([]object.Account, error) {
	a, _ := r.FindByUsername(ctx, username)
	var res []object.Account
	for _, v := range r.db.relationship {
		if v.userID == a.ID {
			ta := r.db.account[v.targetID]
			res = append(res, object.Account{ID: ta.id, Username: ta.username})
		}
	}
	return res, nil
}

func (r *accountMock) Followers(ctx context.Context, username string) ([]object.Account, error) {
	a, _ := r.FindByUsername(ctx, username)
	var res []object.Account
	for _, v := range r.db.relationship {
		if v.targetID == a.ID {
			ua := r.db.account[v.userID]
			res = append(res, object.Account{ID: ua.id, Username: ua.username})
		}
	}
	return res, nil
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
	return &object.Status{
		ID: s.id,
		Account: &object.Account{
			Username: s.username},
	}, nil
}

func (r *statusMock) Create(ctx context.Context, status *object.Status) (object.AccountID, error) {
	newID := int64(len(r.db.status) + 1)
	r.db.status[newID] = statusData{id: newID, username: status.Account.Username}
	return newID, nil
}

func (r *statusMock) Delete(ctx context.Context, status_id object.StatusID, account_id object.AccountID) error {
	delete(r.db.status, status_id)
	return nil
}

func (r *statusMock) AllStatuses(ctx context.Context) ([]object.Status, error) {
	var statuses []object.Status
	for _, value := range r.db.status {
		statuses = append(statuses,
			object.Status{
				ID:      value.id,
				Account: &object.Account{Username: value.username}})
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
	for _, value := range r.db.relationship {
		if value.userID == userID && value.targetID == targetID {
			return true, nil
		}
	}
	return false, nil
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
	newID := int64(len(r.db.status) + 1)
	r.db.relationship[newID] = relationshipData{
		userID:   userID,
		targetID: targetID,
	}
	return newID, nil
}

func (r *relationshipMock) Delete(ctx context.Context, userID object.AccountID, targetID object.AccountID) error {
	return nil
}
