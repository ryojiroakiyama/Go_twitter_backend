package handlertest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"yatter-backend-go/app/app"
)

// ここ引数多いしあんまり使い勝手良くない
// そもそもC structがいるかという, 作るならもっと活用できそう
func Setup(t *testing.T, db *DBMock, newRouter func(app *app.App) http.Handler) *C {
	db = fillDB(db)
	app := &app.App{Dao: &DaoMock{db: db}}
	server := httptest.NewServer(newRouter(app))
	return &C{
		App:    app,
		Server: server,
	}
}

func fillDB(db *DBMock) *DBMock {
	if db == nil {
		db = new(DBMock)
	}
	if db.Account == nil {
		db.Account = make(AccountTableMock)
	}
	if db.Status == nil {
		db.Status = make(StatusTableMock)
	}
	if db.RelationShip == nil {
		db.RelationShip = make(RelationShipTableMock)
	}
	return db
}

type C struct {
	App    *app.App
	Server *httptest.Server
}

func (c *C) Close() {
	c.Server.Close()
}

func (c *C) PostJSON(apiPath string, payload string) (*http.Response, error) {
	return c.Server.Client().Post(c.asURL(apiPath), "application/json", bytes.NewReader([]byte(payload)))
}

//func (c *C) Do(method, apiPath string, body io.Reader, userAuth string) (*http.Response, error) {
//	req, err := http.NewRequest(method, c.asURL(apiPath), body)
//	if err != nil {
//		return nil, err
//	}
//	if body != nil {
//		req.Header.Set("Content-Type", "application/json")
//	}
//	if userAuth != "" {
//		req.Header.Set("Authentication", "username "+userAuth)
//	}
//	return c.Server.Client().Do(req)
//}

func (c *C) Get(apiPath string) (*http.Response, error) {
	return c.Server.Client().Get(c.asURL(apiPath))
}

func (c *C) PostJsonWithAuth(apiPath string, payload string, authUser string) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.asURL(apiPath), bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authentication", "username "+authUser)
	return c.Server.Client().Do(req)
}

func (c *C) asURL(apiPath string) string {
	baseURL, _ := url.Parse(c.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	return baseURL.String()
}
