package handlertest

// Package handlertest はhandler配下のパッケージのテストで使われる機能をまとめたものです。

// NOTE: dbをmapでモックするのをやってみた
// NOTE: テスト時以外でもコンパイルされることが懸念点

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"yatter-backend-go/app/app"
)

// TODO: 機能の整理
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

func (c *C) Get(apiPath string) (*http.Response, error) {
	return c.Server.Client().Get(c.asURL(apiPath))
}

func (c *C) GetWithParam(apiPath string, param string) (*http.Response, error) {
	return c.Server.Client().Get(c.asURL(apiPath) + param)
}

func (c *C) GetWithAuth(apiPath string, authUser string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.asURL(apiPath), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authentication", "username "+authUser)
	return c.Server.Client().Do(req)
}

func (c *C) GetWithParamAuth(apiPath string, param string, authUser string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.asURL(apiPath)+param, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authentication", "username "+authUser)
	return c.Server.Client().Do(req)
}

func (c *C) PostJSON(apiPath string, payload string) (*http.Response, error) {
	return c.Server.Client().Post(c.asURL(apiPath), "application/json", bytes.NewReader([]byte(payload)))
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

func (c *C) DeleteWithAuth(apiPath string, authUser string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.asURL(apiPath), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authentication", "username "+authUser)
	return c.Server.Client().Do(req)
}

func (c *C) asURL(apiPath string) string {
	baseURL, _ := url.Parse(c.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	return baseURL.String()
}
