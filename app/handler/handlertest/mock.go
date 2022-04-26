package handlertest

import (
	"context"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"
)

/*
 * dao, db, repojitoryをモック
 * 構成はdaoと同じようにした
 */

// DaoMock: daoのモック
type DaoMock struct {
	db *DBMock
}

func (d *DaoMock) Account() repository.Account {
	return newAccountMock(d.db)
}

func (d *DaoMock) Status() repository.Status {
	return newStatusMock(d.db)
}

func (d *DaoMock) Relationship() repository.Relationship {
	return newRelationShipMock(d.db)
}

func (d *DaoMock) InitAll() error {
	d.db = nil
	return nil
}

// DBMock: 各テーブルをmapで模したdbモック
type DBMock struct {
	Account      AccountTableMock
	Status       StatusTableMock
	RelationShip RelationShipTableMock
}

type AccountTableMock map[int64]AccountData
type AccountData struct {
	ID       int64
	UserName string
}

type StatusTableMock map[int64]StatusData
type StatusData struct {
	ID       int64
	content  string
	UserName string
}

type RelationShipTableMock map[int64]RelationShipData
type RelationShipData struct {
	userID   int64
	targetID int64
}

// accountMock: Account repojitoryをモック
type accountMock struct {
	db *DBMock
}

func newAccountMock(db *DBMock) repository.Account {
	return &accountMock{db: db}
}

func (r *accountMock) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	for _, value := range r.db.Account {
		if value.UserName == username {
			return &object.Account{
					ID:       value.ID,
					Username: value.UserName},
				nil
		}
	}
	return nil, nil
}

func (r *accountMock) Create(ctx context.Context, account *object.Account) (object.AccountID, error) {
	newID := int64(len(r.db.Account) + 1)
	r.db.Account[int64(newID)] = AccountData{
		ID:       newID,
		UserName: account.Username}
	return newID, nil
}

func (r *accountMock) Following(ctx context.Context, username string) ([]object.Account, error) {
	a, _ := r.FindByUsername(ctx, username)
	var res []object.Account
	for _, v := range r.db.RelationShip {
		if v.userID == a.ID {
			ta := r.db.Account[v.targetID]
			res = append(res, object.Account{ID: ta.ID, Username: ta.UserName})
		}
	}
	return res, nil
}

func (r *accountMock) Followers(ctx context.Context, username string) ([]object.Account, error) {
	a, _ := r.FindByUsername(ctx, username)
	var res []object.Account
	for _, v := range r.db.RelationShip {
		if v.targetID == a.ID {
			ua := r.db.Account[v.userID]
			res = append(res, object.Account{ID: ua.ID, Username: ua.UserName})
		}
	}
	return res, nil
}

// statusMock: Status repojitoryをモック
type statusMock struct {
	db *DBMock
}

func newStatusMock(db *DBMock) repository.Status {
	return &statusMock{db: db}
}

func (r *statusMock) FindByID(ctx context.Context, id object.StatusID) (*object.Status, error) {
	s, exist := r.db.Status[id]
	if !exist {
		return nil, nil
	}
	return &object.Status{
		ID:      s.ID,
		Content: s.content,
		Account: &object.Account{
			Username: s.UserName},
	}, nil
}

func (r *statusMock) Create(ctx context.Context, Status *object.Status) (object.AccountID, error) {
	newID := int64(len(r.db.Status) + 1)
	r.db.Status[newID] = StatusData{ID: newID, UserName: Status.Account.Username}
	return newID, nil
}

func (r *statusMock) Delete(ctx context.Context, status_id object.StatusID, account_id object.AccountID) error {
	delete(r.db.Status, status_id)
	return nil
}

func (r *statusMock) AllStatuses(ctx context.Context) ([]object.Status, error) {
	var statuses []object.Status
	for _, v := range r.db.Status {
		statuses = append(statuses,
			object.Status{
				ID:      v.ID,
				Content: v.content,
				Account: &object.Account{Username: v.UserName}})
	}
	return statuses, nil
}

func (r *statusMock) FollowingStatuses(ctx context.Context, username string) ([]object.Status, error) {
	return nil, nil
}

// relationshipMock: RelationShip repojitoryをモック
type relationshipMock struct {
	db *DBMock
}

func newRelationShipMock(db *DBMock) repository.Relationship {
	return &relationshipMock{db: db}
}

func (r *relationshipMock) IsFollowing(ctx context.Context, userID object.AccountID, targetID object.AccountID) (bool, error) {
	for _, value := range r.db.RelationShip {
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
	newID := int64(len(r.db.Status) + 1)
	r.db.RelationShip[newID] = RelationShipData{
		userID:   userID,
		targetID: targetID,
	}
	return newID, nil
}

func (r *relationshipMock) Delete(ctx context.Context, userID object.AccountID, targetID object.AccountID) error {
	for ID, v := range r.db.RelationShip {
		if v.userID == userID && v.targetID == targetID {
			delete(r.db.RelationShip, ID)
		}
	}
	return nil
}
