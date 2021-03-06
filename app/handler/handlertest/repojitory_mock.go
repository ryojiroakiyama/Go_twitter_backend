package handlertest

import (
	"context"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"
)

// accountMock: Account repojitoryを実装
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

func (r *accountMock) Update(ctx context.Context, entity *object.Account) error {
	return nil
}

func (r *accountMock) Followings(ctx context.Context, username string, limit int64) ([]object.Account, error) {
	a, _ := r.FindByUsername(ctx, username)
	var res []object.Account
	for _, v := range r.db.RelationShip {
		if v.UserID == a.ID {
			ta := r.db.Account[v.TargetID]
			res = append(res, object.Account{ID: ta.ID, Username: ta.UserName})
		}
	}
	sortAccounts(res)
	if int64(len(res)) < limit {
		return res, nil
	} else {
		return res[:limit], nil
	}
}

func (r *accountMock) Followers(ctx context.Context, username string, since_id int64, max_id int64, limit int64) ([]object.Account, error) {
	a, _ := r.FindByUsername(ctx, username)
	var res []object.Account
	for _, v := range r.db.RelationShip {
		if v.TargetID == a.ID && since_id <= v.UserID && v.UserID <= max_id {
			ua := r.db.Account[v.UserID]
			res = append(res, object.Account{ID: ua.ID, Username: ua.UserName})
		}
	}
	sortAccounts(res)
	if int64(len(res)) < limit {
		return res, nil
	} else {
		return res[:limit], nil
	}
}

// statusMock: Status repojitoryを実装
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
		Content: s.Content,
		Account: &object.Account{
			Username: s.UserName},
	}, nil
}

func (r *statusMock) Create(ctx context.Context, Status *object.Status) (object.AccountID, error) {
	newID := int64(len(r.db.Status) + 1)
	r.db.Status[newID] = StatusData{ID: newID, UserName: Status.Account.Username, Content: Status.Content}
	return newID, nil
}

func (r *statusMock) Delete(ctx context.Context, status_id object.StatusID, account_id object.AccountID) error {
	delete(r.db.Status, status_id)
	return nil
}

// TODO: parameter対応
func (r *statusMock) AllStatuses(ctx context.Context, since_id int64, max_id int64, limit int64) ([]object.Status, error) {
	var statuses []object.Status
	for _, v := range r.db.Status {
		statuses = append(statuses,
			object.Status{
				ID:      v.ID,
				Content: v.Content,
				Account: &object.Account{Username: v.UserName}})
	}
	sortStatuses(statuses)
	return statuses, nil
}

// 今からこれ再現するのはコスト高いので保留
func (r *statusMock) RelationStatuses(ctx context.Context, user_id object.AccountID, since_id int64, max_id int64, limit int64) ([]object.Status, error) {
	return []object.Status{}, nil
}

// relationshipMock: RelationShip repojitoryを実装
type relationshipMock struct {
	db *DBMock
}

func newRelationShipMock(db *DBMock) repository.Relationship {
	return &relationshipMock{db: db}
}

func (r *relationshipMock) IsFollowing(ctx context.Context, userID object.AccountID, targetID object.AccountID) (bool, error) {
	for _, value := range r.db.RelationShip {
		if value.UserID == userID && value.TargetID == targetID {
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
		UserID:   userID,
		TargetID: targetID,
	}
	return newID, nil
}

func (r *relationshipMock) Delete(ctx context.Context, userID object.AccountID, targetID object.AccountID) error {
	for ID, v := range r.db.RelationShip {
		if v.UserID == userID && v.TargetID == targetID {
			delete(r.db.RelationShip, ID)
		}
	}
	return nil
}

// mapで取り出されたarrayは順番がランダムなので用意した
func sortAccounts(a []object.Account) {
	l := len(a)
	for i := int64(0); i < int64(l); i++ {
		for j := i + 1; j < int64(l); j++ {
			if a[i].ID > a[j].ID {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
}

func sortStatuses(a []object.Status) {
	l := len(a)
	for i := int64(0); i < int64(l); i++ {
		for j := i + 1; j < int64(l); j++ {
			if a[i].ID > a[j].ID {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
}
