package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"

	"yatter-backend-go/app/app"
	"yatter-backend-go/app/domain/object"

	"github.com/stretchr/testify/assert"
)

func TestAccountRegistration(t *testing.T) {
	c := setup(t)
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

func setup(t *testing.T) *C {
	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	if err := app.Dao.InitAll(); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(NewRouter(app))

	return &C{
		App:    app,
		Server: server,
	}
}

// mapは重複してたら上書きしてしまうので, insert前に確認忘れずに
type account struct {
	account map[string]*object.Account
}

func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	account, exist := r.account[username]
	if exist {
		return account, nil
	}
	return nil, fmt.Errorf("FindByUsername: not exist")
}

// Create: アカウント作成
func (r *account) Create(ctx context.Context, entity *object.Account) (object.AccountID, error) {
	_, exist := r.account[entity.Username]
	if exist {
		return 0, fmt.Errorf("Create: aready exist")
	}
	id := len(r.account) + 1
	r.account[entity.Username] = entity
	return int64(id), nil
}

type status struct {
	status map[object.StatusID]*object.Status
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
