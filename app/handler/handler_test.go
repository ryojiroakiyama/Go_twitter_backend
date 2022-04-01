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
)

func TestAccountRegistration(t *testing.T) {
	john := &object.Account{
		Username: "john",
	}
	status1 := &object.Status{
		ID:      1,
		Content: "status1",
	}
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
			name:           "account create normal",
			method:         "POST",
			apiPath:        "/v1/accounts",
			body:           bytes.NewReader([]byte(`{"username":"john"}`)),
			bodyExpected:   john,
			statusExpected: http.StatusOK,
		},
		//{
		//	name: "account create duplicate",
		//	db: func() *dbMock {
		//		a := make(accountTableMock)
		//		s := make(statusTableMock)
		//		a[john.Username] = john
		//		return &dbMock{account: a, status: s}
		//	}(),
		//	method:         "POST",
		//	apiPath:        "/v1/accounts",
		//	body:           bytes.NewReader([]byte(`{"username":"john"}`)),
		//	bodyExpected:   john,
		//	statusExpected: http.StatusOK,
		//},
		{
			name: "account fetch",
			db: func() *dbMock {
				a := make(accountTableMock)
				a[john.Username] = john
				return &dbMock{account: a}
			}(),
			method:         "GET",
			apiPath:        "/v1/accounts/john",
			bodyExpected:   john,
			statusExpected: http.StatusOK,
		},
		{
			name: "status fetch",
			db: func() *dbMock {
				s := make(statusTableMock)
				s[1] = status1
				return &dbMock{status: s}
			}(),
			method:         "GET",
			apiPath:        "/v1/statuses/1",
			bodyExpected:   status1,
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
			defer resp.Body.Close()
			if resp.StatusCode != tt.statusExpected {
				t.Fatalf("expected: %v, returned: %v", tt.statusExpected, resp.StatusCode)
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
}

func setup(t *testing.T, db *dbMock) *C {
	if db == nil {
		db = new(dbMock)
	}
	if db.account == nil {
		db.account = make(accountTableMock)
	}
	if db.status == nil {
		db.status = make(statusTableMock)
	}

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
	a, exist := r.db.account[username]
	if !exist {
		return nil, fmt.Errorf("FindByUsername: Account not exist")
	}
	return a, nil
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
	s, exist := r.db.status[id]
	if !exist {
		return nil, fmt.Errorf("FindByID: Status not exist")
	}
	return s, nil
}

func (r *statusMock) Create(ctx context.Context, entity *object.Status) (object.AccountID, error) {
	entity.ID = int64(len(r.db.status) + 1)
	r.db.status[entity.ID] = entity
	return entity.ID, nil
}

func (r *statusMock) Delete(ctx context.Context, status_id object.StatusID, account_id object.AccountID) error {
	s, exist := r.db.status[status_id]
	if !exist {
		return fmt.Errorf("Delete: No status matched")
	} else if s.Account.ID != account_id {
		return fmt.Errorf("Delete: No status matched")
	}
	delete(r.db.status, status_id)
	return nil
}

func (r *statusMock) All(ctx context.Context) ([]object.Status, error) {
	var statuses []object.Status
	for _, value := range r.db.status {
		statuses = append(statuses, *value)
	}
	return statuses, nil
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

func (d *daoMock) Relationship() repository.Relationship {
	return nil
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
	req.Header.Set("Content-Type", "application/json")
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
