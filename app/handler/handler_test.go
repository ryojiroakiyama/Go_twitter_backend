package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"yatter-backend-go/app/domain/object"
)

type TestSource []struct {
	Name         string
	DB           *dbMock
	Method       string
	ApiPath      string
	Body         io.Reader
	AuthUserName string
	WantBody     []byte
	WantStatus   int
}

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
	tests := TestSource{
		{
			Name:       "create account",
			Method:     "POST",
			ApiPath:    "/v1/accounts",
			Body:       bytes.NewReader([]byte(`{"username":"john"}`)),
			WantBody:   ToJsonFormat(t, john),
			WantStatus: http.StatusOK,
		},
		{
			Name: "create duplicate account",
			DB: func() *dbMock {
				a := make(accountTableMock)
				a[john.Username] = *john
				return &dbMock{account: a}
			}(),
			Method:     "POST",
			ApiPath:    "/v1/accounts",
			Body:       bytes.NewReader([]byte(`{"username":"john"}`)),
			WantBody:   []byte(http.StatusText(http.StatusConflict) + "\n"),
			WantStatus: http.StatusConflict,
		},
		{
			Name: "fetch account",
			DB: func() *dbMock {
				a := make(accountTableMock)
				a[john.Username] = *john
				return &dbMock{account: a}
			}(),
			Method:     "GET",
			ApiPath:    "/v1/accounts/john",
			WantBody:   ToJsonFormat(t, john),
			WantStatus: http.StatusOK,
		},
		{
			Name:       "fetch non-exist account",
			Method:     "GET",
			ApiPath:    "/v1/accounts/john",
			WantBody:   []byte(http.StatusText(http.StatusNotFound) + "\n"),
			WantStatus: http.StatusNotFound,
		},
		{
			Name: "fetch status",
			DB: func() *dbMock {
				s := make(statusTableMock)
				s[1] = *johnStatus1
				return &dbMock{status: s}
			}(),
			Method:     "GET",
			ApiPath:    "/v1/statuses/1",
			WantBody:   ToJsonFormat(t, johnStatus1),
			WantStatus: http.StatusOK,
		},
		{
			Name:       "fetch non-exist status",
			Method:     "GET",
			ApiPath:    "/v1/statuses/1",
			WantBody:   []byte(http.StatusText(http.StatusNotFound) + "\n"),
			WantStatus: http.StatusNotFound,
		},
		{
			Name: "create status",
			DB: func() *dbMock {
				a := make(accountTableMock)
				a[john.Username] = *john
				return &dbMock{account: a}
			}(),
			Method:       "POST",
			ApiPath:      "/v1/statuses",
			Body:         bytes.NewReader([]byte(`{"status":"johnStatus"}`)),
			AuthUserName: john.Username,
			WantBody:     ToJsonFormat(t, johnStatus1),
			WantStatus:   http.StatusOK,
		},
		{
			Name: "create duplicate status",
			DB: func() *dbMock {
				a := make(accountTableMock)
				s := make(statusTableMock)
				a[john.Username] = *john
				s[1] = *johnStatus1
				return &dbMock{account: a, status: s}
			}(),
			Method:       "POST",
			ApiPath:      "/v1/statuses",
			Body:         bytes.NewReader([]byte(`{"status":"johnStatus"}`)),
			AuthUserName: john.Username,
			WantBody:     ToJsonFormat(t, johnStatus2),
			WantStatus:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			c := setup(t, tt.DB)
			defer c.Close()

			resp, err := c.Do(tt.Method, tt.ApiPath, tt.Body, tt.AuthUserName)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// check status code
			if resp.StatusCode != tt.WantStatus {
				t.Fatalf("code expected: %v, returned: %v", tt.WantStatus, resp.StatusCode)
			}

			// check response body
			if tt.WantBody != nil {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				if bytes.Compare(body, tt.WantBody) != 0 {
					t.Fatalf("body \nexpected: [%v], \nreturned: [%v]", string(tt.WantBody), string(body))
				}
			}
		})
	}
}

func ToJsonFormat(t *testing.T, body interface{}) []byte {
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
