package accounts_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"yatter-backend-go/app/domain/object"
)

func TestCreate(t *testing.T) {
	john := &object.Account{
		Username: "john",
	}
	tests := []struct {
		name         string
		db           *dbMock
		method       string
		apiPath      string
		body         io.Reader
		authUserName string
		wantBody     []byte
		wantStatus   int
	}{
		{
			name:       "create account",
			method:     "POST",
			apiPath:    "/",
			body:       bytes.NewReader([]byte(`{"username":"john"}`)),
			wantBody:   toJsonFormat(t, john),
			wantStatus: http.StatusOK,
		},
		{
			name: "create duplicate account",
			db: func() *dbMock {
				a := make(accountTableMock)
				a[john.Username] = *john
				return &dbMock{account: a}
			}(),
			method:     "POST",
			apiPath:    "/",
			body:       bytes.NewReader([]byte(`{"username":"john"}`)),
			wantBody:   nil, // bodyの確認省略
			wantStatus: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := setup(t, tt.db)
			defer c.Close()

			resp, err := c.Do(tt.method, tt.apiPath, tt.body, tt.authUserName)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// check status code
			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("code expected: %v, returned: %v", tt.wantStatus, resp.StatusCode)
			}

			// check response body
			if tt.wantBody != nil {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				if bytes.Compare(body, tt.wantBody) != 0 {
					t.Fatalf("body \nexpected: [%v], \nreturned: [%v]", string(tt.wantBody), string(body))
				}
			}
		})
	}
}
