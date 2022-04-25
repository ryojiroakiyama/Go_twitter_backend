package accounts_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler/accounts"
)

func setup(t *testing.T, db *dbMock) *C {
	db = fillDB(db)
	app := &app.App{Dao: &daoMock{db: db}}
	server := httptest.NewServer(accounts.NewRouter(app))
	return &C{
		App:    app,
		Server: server,
	}
}

func fillDB(db *dbMock) *dbMock {
	if db == nil {
		db = new(dbMock)
	}
	if db.account == nil {
		db.account = make(accountTableMock)
	}
	if db.status == nil {
		db.status = make(statusTableMock)
	}
	if db.relationship == nil {
		db.relationship = make(relationshipTableMock)
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

func (c *C) Do(method, apiPath string, body io.Reader, userAuth string) (*http.Response, error) {
	req, err := http.NewRequest(method, c.asURL(apiPath), body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if userAuth != "" {
		req.Header.Set("Authentication", "username "+userAuth)
	}
	return c.Server.Client().Do(req)
}

func (c *C) Get(apiPath string) (*http.Response, error) {
	return c.Server.Client().Get(c.asURL(apiPath))
}

func (c *C) PostWithAuth(apiPath string, payload string, authUser string) (*http.Response, error) {
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
