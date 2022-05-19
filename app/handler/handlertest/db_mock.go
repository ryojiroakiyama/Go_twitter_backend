package handlertest

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
	Content  string
	UserName string
}

type RelationShipTableMock map[int64]RelationShipData
type RelationShipData struct {
	UserID   int64
	TargetID int64
}

type MediaTableMock map[int64]MediaData
type MediaData struct {
	ID  int64
	URL string
}
