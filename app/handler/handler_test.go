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
	"yatter-backend-go/app/handler/accounts"
)

func TestAccountRegistration(t *testing.T) {
	john := &object.Account{
		Username: "john",
	}
	johnStatus := &object.Status{
		ID:      1,
		Content: "johnStatus",
		Account: john,
	}
	tests := []struct {
		name           string
		db             *dbMock
		method         string
		apiPath        string
		body           io.Reader
		userAuth       string
		bodyExpected   []byte
		statusExpected int
	}{
		{
			name:           "create account",
			method:         "POST",
			apiPath:        "/v1/accounts",
			body:           bytes.NewReader([]byte(`{"username":"john"}`)),
			bodyExpected:   jsonFormat(t, john),
			statusExpected: http.StatusOK,
		},
		{
			name: "create account duplicate",
			db: func() *dbMock {
				a := make(accountTableMock)
				a[john.Username] = *john
				return &dbMock{account: a}
			}(),
			method:         "POST",
			apiPath:        "/v1/accounts",
			body:           bytes.NewReader([]byte(`{"username":"john"}`)),
			bodyExpected:   []byte(accounts.TextUserConflict + "\n"),
			statusExpected: http.StatusConflict,
		},
		{
			name: "fetch account",
			db: func() *dbMock {
				a := make(accountTableMock)
				a[john.Username] = *john
				return &dbMock{account: a}
			}(),
			method:         "GET",
			apiPath:        "/v1/accounts/john",
			bodyExpected:   jsonFormat(t, john),
			statusExpected: http.StatusOK,
		},
		{
			name:           "fetch no exist account",
			db:             nil,
			method:         "GET",
			apiPath:        "/v1/accounts/john",
			bodyExpected:   []byte(accounts.TextNoAccount + "\n"),
			statusExpected: http.StatusNotFound,
		},
		{
			name: "fetch status",
			db: func() *dbMock {
				s := make(statusTableMock)
				s[1] = *johnStatus
				return &dbMock{status: s}
			}(),
			method:         "GET",
			apiPath:        "/v1/statuses/1",
			bodyExpected:   jsonFormat(t, johnStatus),
			statusExpected: http.StatusOK,
		},
		{
			name: "create status",
			db: func() *dbMock {
				a := make(accountTableMock)
				a[john.Username] = *john
				return &dbMock{account: a}
			}(),
			method:         "POST",
			apiPath:        "/v1/statuses",
			body:           bytes.NewReader([]byte(`{"status":"johnStatus"}`)),
			userAuth:       john.Username,
			bodyExpected:   jsonFormat(t, johnStatus),
			statusExpected: http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.db = fillDB(tt.db)
			c := setup(t, tt.db)
			defer c.Close()

			resp, err := c.Do(tt.method, tt.apiPath, tt.body, tt.userAuth)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// check status code
			if resp.StatusCode != tt.statusExpected {
				t.Fatalf("code expected: %v, returned: %v", tt.statusExpected, resp.StatusCode)
			}

			// check response body
			if tt.bodyExpected != nil {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				if bytes.Compare(body, tt.bodyExpected) != 0 {
					t.Fatalf("body \nexpected: [%v], \nreturned: [%v]", string(tt.bodyExpected), string(body))
				}
			}
		})
	}
}

func jsonFormat(t *testing.T, body interface{}) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}
	return out
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
	return db
}

func setup(t *testing.T, db *dbMock) *C {
	app := &app.App{Dao: &daoMock{db: db}}
	server := httptest.NewServer(NewRouter(app))
	return &C{
		App:    app,
		Server: server,
	}
}

type accountTableMock map[string]object.Account
type statusTableMock map[object.StatusID]object.Status

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
		return nil, nil
	}
	return &a, nil
}

func (r *accountMock) Create(ctx context.Context, entity *object.Account) (object.AccountID, error) {
	_, exist := r.db.account[entity.Username]
	if exist {
		return 0, fmt.Errorf("Create: Account aready exist")
	}
	r.db.account[entity.Username] = *entity
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
	return &s, nil
}

func (r *statusMock) Create(ctx context.Context, entity *object.Status) (object.AccountID, error) {
	entity.ID = int64(len(r.db.status) + 1)
	r.db.status[entity.ID] = *entity
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
		statuses = append(statuses, value)
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

func (c *C) asURL(apiPath string) string {
	baseURL, _ := url.Parse(c.Server.URL)
	baseURL.Path = path.Join(baseURL.Path, apiPath)
	return baseURL.String()
}
