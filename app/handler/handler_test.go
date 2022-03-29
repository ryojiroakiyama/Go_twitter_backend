package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"yatter-backend-go/app/app"
	"yatter-backend-go/app/dao"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/stretchr/testify/assert"
)

func TestAccountRegistration(t *testing.T) {
	expect := &mockDao{}
	expect.account.findbyusername.obj = &object.Account{
		Username: "john",
	}
	c := setup(t, expect)
	defer c.Close()

	func() {
		resp, err := c.PostJSON("/v1/accounts", `{"username":"john"}`)
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, resp.StatusCode, http.StatusOK) {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, "john", j["username"])
		}
	}()

	func() {
		resp, err := c.Get("/v1/accounts/john")
		if err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, resp.StatusCode, http.StatusOK) {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var j map[string]interface{}
		if assert.NoError(t, json.Unmarshal(body, &j)) {
			assert.Equal(t, "john", j["username"])
		}
	}()
}

func setup(t *testing.T, d dao.Dao) *C {
	app := &app.App{Dao: d}

	if err := app.Dao.InitAll(); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(NewRouter(app))

	return &C{
		App:    app,
		Server: server,
	}
}

type retFindBy struct {
	obj interface{}
	err error
}

type retCreate struct {
	id  int64
	err error
}

type retDelete struct {
	err error
}

type retAll struct {
	obj interface{}
	err error
}

type mockAccount struct {
	findbyusername retFindBy
	create         retCreate
}

func (r *mockAccount) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	obj, ok := r.findbyusername.obj.(*object.Account)
	if !ok {
		panic("findbyuser")
	}
	return obj, r.findbyusername.err
}

func (r *mockAccount) Create(ctx context.Context, entity *object.Account) (object.AccountID, error) {
	return r.create.id, r.create.err
}

type mockStatus struct {
	findbyid retFindBy
	create   retCreate
	delete   retDelete
	all      retAll
}

func (r *mockStatus) FindByID(ctx context.Context, id object.StatusID) (*object.Status, error) {
	obj, ok := r.findbyid.obj.(*object.Status)
	if !ok {
		panic("findbyid")
	}
	return obj, r.findbyid.err
}

func (r *mockStatus) Create(ctx context.Context, entity *object.Status) (object.AccountID, error) {
	return r.create.id, r.create.err
}

func (r *mockStatus) Delete(ctx context.Context, status_id object.StatusID, account_id object.AccountID) error {
	return r.delete.err
}

func (r *mockStatus) All(ctx context.Context) ([]object.Status, error) {
	obj, ok := r.all.obj.([]object.Status)
	if !ok {
		panic("all")
	}
	return obj, r.all.err
}

type mockDao struct {
	account mockAccount
	status  mockStatus
}

func (d *mockDao) Account() repository.Account {
	return &d.account
}

func (d *mockDao) Status() repository.Status {
	return &d.status
}

func (d *mockDao) InitAll() error {
	return nil
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

func (c *C) Get(apiPath string) (*http.Response, error) {
	return c.Server.Client().Get(c.asURL(apiPath))
}

func (c *C) asURL(apiPath string) string {
	baseURL, _ := url.Parse(c.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	return baseURL.String()
}
