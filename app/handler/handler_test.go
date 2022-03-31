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
	"yatter-backend-go/app/domain/repository"

	"github.com/stretchr/testify/assert"
)

func TestAccountRegistration(t *testing.T) {
	tests := []struct {
		name           string
		accountDB      map[string]*object.Account
		statusDB       map[object.StatusID]*object.Status
		method         string
		apiPath        string
		body           io.Reader
		bodyExpected   interface{}
		statusExpected int
	}{
		{
			name: "account fetch",
			accountDB: map[string]*object.Account{
				"john": {
					Username: "john",
				},
			},
			method:  "GET",
			apiPath: "/v1/accounts/john",
			bodyExpected: &object.Account{
				Username: "john",
			},
			statusExpected: http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var a accountMock
			var s statusMock
			a.dbMock = tt.accountDB
			s.dbMock = tt.statusDB
			c := setup(t, &a, &s)
			defer c.Close()

			resp, err := c.Do(tt.method, tt.apiPath, tt.body)
			if err != nil {
				t.Fatal(err)
			}
			if !assert.Equal(t, resp.StatusCode, tt.statusExpected) {
				return
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			buf := new(bytes.Buffer)
			if err := json.NewEncoder(buf).Encode(tt.bodyExpected); err != nil {
				t.Fatal(err)
			}
			expect, err := io.ReadAll(buf)
			if bytes.Compare(body, expect) != 0 {
				t.Fatalf("expected: %v, returned: %v", string(expect), string(body))
			}
		})
	}
	//expect := &daoMock{}
	//expect.account.findbyusername.obj = &object.Account{
	//	Username: "john",
	//}
	//c := setup(t, expect)
	//defer c.Close()

	//func() {
	//	resp, err := c.PostJSON("/v1/accounts", `{"username":"john"}`)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	if !assert.Equal(t, resp.StatusCode, http.StatusOK) {
	//		return
	//	}

	//	body, err := io.ReadAll(resp.Body)
	//	if err != nil {
	//		t.Fatal(err)
	//	}

	//	var j map[string]interface{}
	//	if assert.NoError(t, json.Unmarshal(body, &j)) {
	//		assert.Equal(t, "john", j["username"])
	//	}
	//}()

	//func() {
	//	resp, err := c.Get("/v1/accounts/john")
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	if !assert.Equal(t, resp.StatusCode, http.StatusOK) {
	//		return
	//	}

	//	body, err := io.ReadAll(resp.Body)
	//	if err != nil {
	//		t.Fatal(err)
	//	}

	//	var j map[string]interface{}
	//	if assert.NoError(t, json.Unmarshal(body, &j)) {
	//		assert.Equal(t, "john", j["username"])
	//	}
	//}()
}

func setup(t *testing.T, a repository.Account, s repository.Status) *C {
	app := &app.App{Dao: &daoMock{account: a, status: s}}

	if err := app.Dao.InitAll(); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(NewRouter(app))

	return &C{
		App:    app,
		Server: server,
	}
}

// accountMock: accountに関するrepojitoryとdbをモック
type accountMock struct {
	dbMock map[string]*object.Account
}

func (r *accountMock) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	accountMock, exist := r.dbMock[username]
	if exist {
		return accountMock, nil
	}
	return nil, fmt.Errorf("FindByUsername: Account not exist")
}

func (r *accountMock) Create(ctx context.Context, entity *object.Account) (object.AccountID, error) {
	_, exist := r.dbMock[entity.Username]
	if exist {
		return 0, fmt.Errorf("Create: Account aready exist")
	}
	id := len(r.dbMock) + 1
	r.dbMock[entity.Username] = entity
	return int64(id), nil
}

// statusMock: statusに関するrepojitoryとdbをモック
type statusMock struct {
	dbMock map[object.StatusID]*object.Status
}

func (r *statusMock) FindByID(ctx context.Context, id object.StatusID) (*object.Status, error) {
	return nil, nil
}

func (r *statusMock) Create(ctx context.Context, entity *object.Status) (object.AccountID, error) {
	return 0, nil
}

func (r *statusMock) Delete(ctx context.Context, status_id object.StatusID, account_id object.AccountID) error {
	return nil
}

func (r *statusMock) All(ctx context.Context) ([]object.Status, error) {
	return nil, nil
}

type daoMock struct {
	account repository.Account
	status  repository.Status
}

func (d *daoMock) Account() repository.Account {
	return d.account
}

func (d *daoMock) Status() repository.Status {
	return d.status
}

func (d *daoMock) InitAll() error {
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

func (c *C) Do(method, apiPath string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.asURL(apiPath), body)
	if err != nil {
		return nil, err
	}
	return c.Server.Client().Do(req)
}

func (c *C) Get(apiPath string) (*http.Response, error) {
	return c.Server.Client().Get(c.asURL(apiPath))
}

func (c *C) asURL(apiPath string) string {
	baseURL, _ := url.Parse(c.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	return baseURL.String()
}
