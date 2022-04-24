package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"yatter-backend-go/app/domain/object"
)

func TestAccountRegistration(t *testing.T) {
	john := &object.Account{
		Username: "john",
	}
	johnStatus1 := &object.Status{
		ID:      1,
		Content: "johnStatus",
		Account: john,
	}
	johnStatus2 := &object.Status{
		ID:      2,
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
			name: "create duplicate account",
			db: func() *dbMock {
				a := make(accountTableMock)
				a[john.Username] = *john
				return &dbMock{account: a}
			}(),
			method:         "POST",
			apiPath:        "/v1/accounts",
			body:           bytes.NewReader([]byte(`{"username":"john"}`)),
			bodyExpected:   []byte(http.StatusText(http.StatusConflict) + "\n"),
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
			name:           "fetch non-exist account",
			method:         "GET",
			apiPath:        "/v1/accounts/john",
			bodyExpected:   []byte(http.StatusText(http.StatusNotFound) + "\n"),
			statusExpected: http.StatusNotFound,
		},
		{
			name: "fetch status",
			db: func() *dbMock {
				s := make(statusTableMock)
				s[1] = *johnStatus1
				return &dbMock{status: s}
			}(),
			method:         "GET",
			apiPath:        "/v1/statuses/1",
			bodyExpected:   jsonFormat(t, johnStatus1),
			statusExpected: http.StatusOK,
		},
		{
			name:           "fetch non-exist status",
			method:         "GET",
			apiPath:        "/v1/statuses/1",
			bodyExpected:   []byte(http.StatusText(http.StatusNotFound) + "\n"),
			statusExpected: http.StatusNotFound,
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
			bodyExpected:   jsonFormat(t, johnStatus1),
			statusExpected: http.StatusOK,
		},
		{
			name: "create duplicate status",
			db: func() *dbMock {
				a := make(accountTableMock)
				s := make(statusTableMock)
				a[john.Username] = *john
				s[1] = *johnStatus1
				return &dbMock{account: a, status: s}
			}(),
			method:         "POST",
			apiPath:        "/v1/statuses",
			body:           bytes.NewReader([]byte(`{"status":"johnStatus"}`)),
			userAuth:       john.Username,
			bodyExpected:   jsonFormat(t, johnStatus2),
			statusExpected: http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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
