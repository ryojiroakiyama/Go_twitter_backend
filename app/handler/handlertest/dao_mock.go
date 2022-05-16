package handlertest

import (
	"yatter-backend-go/app/domain/repository"
)

/*
 * 構成はソースパッケージのdaoと同じようにした
 */

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

func (d *DaoMock) Media() repository.Media {
	return nil
}

func (d *DaoMock) InitAll() error {
	d.db = nil
	return nil
}
