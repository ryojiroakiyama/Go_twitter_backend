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
		name       string
		db         *dbMock
		body       io.Reader
		wantStatus int
		toTestBody bool
		wantBody   []byte
	}{
		{
			name:       "create account",
			body:       bytes.NewReader([]byte(`{"username":"john"}`)),
			wantStatus: http.StatusOK,
			toTestBody: true,
			wantBody:   toJsonFormat(t, john),
		},
		{
			name: "duplicate account",
			db: func() *dbMock {
				a := make(accountTableMock)
				a[john.Username] = *john
				return &dbMock{account: a}
			}(),
			body:       bytes.NewReader([]byte(`{"username":"john"}`)),
			wantStatus: http.StatusConflict,
			toTestBody: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := setup(t, tt.db)
			defer c.Close()

			resp, err := c.Do("POST", "/", tt.body, "")
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// check status code
			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("code expected: %v, returned: %v", tt.wantStatus, resp.StatusCode)
			}

			// check response body
			if tt.toTestBody {
				responseBody, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				if bytes.Compare(responseBody, tt.wantBody) != 0 {
					t.Fatalf("body \nexpected: [%v], \nreturned: [%v]", string(tt.wantBody), string(responseBody))
				}
			}
		})
	}
}
