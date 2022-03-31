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
		db             *dbMock
		method         string
		apiPath        string
		body           io.Reader
		bodyExpected   interface{}
		statusExpected int
	}{
		{
			name: "account fetch",
			db: &dbMock{
				account: accountTableMock{
					"john": {
						Username: "john",
					},
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
			c := setup(t, tt.db)
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

func setup(t *testing.T, db *dbMock) *C {
	app := &app.App{Dao: &daoMock{db: db}}

	server := httptest.NewServer(NewRouter(app))

	return &C{
		App:    app,
		Server: server,
	}
}

type accountTableMock map[string]*object.Account
type statusTableMock map[object.StatusID]*object.Status

type dbMock struct {
	account accountTableMock
	status  statusTableMock
}

// accountMock: accountに関するrepojitoryとtableをモック
type accountMock struct {
	db *dbMock
}

func newAccountMock(db *dbMock) repository.Account {
	return &accountMock{db: db}
}

func (r *accountMock) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	accountMock, exist := r.db.account[username]
	if exist {
		return accountMock, nil
	}
	return nil, fmt.Errorf("FindByUsername: Account not exist")
}

func (r *accountMock) Create(ctx context.Context, entity *object.Account) (object.AccountID, error) {
	_, exist := r.db.account[entity.Username]
	if exist {
		return 0, fmt.Errorf("Create: Account aready exist")
	}
	r.db.account[entity.Username] = entity
	id := len(r.db.account) + 1
	return int64(id), nil
}

// statusMock: statusに関するrepojitoryとtableをモック
type statusMock struct {
	db *dbMock
}

func newStatusMock(db *dbMock) repository.Status {
	return &statusMock{db: db}
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
	db *dbMock
}

func (d *daoMock) Account() repository.Account {
	return newAccountMock(d.db)
}

func (d *daoMock) Status() repository.Status {
	return newStatusMock(d.db)
}

func (d *daoMock) InitAll() error {
	d.db = nil
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
