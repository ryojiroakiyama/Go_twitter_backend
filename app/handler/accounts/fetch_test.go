package accounts_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"yatter-backend-go/app/domain/object"
)

func TestFetch(t *testing.T) {
	john := &object.Account{
		Username: "john",
	}
	tests := []struct {
		name       string
		db         *dbMock
		apiPath    string
		body       io.Reader
		wantBody   []byte
		wantStatus int
	}{
		{
			name: "simple",
			db: func() *dbMock {
				a := make(accountTableMock)
				a[john.Username] = *john
				return &dbMock{account: a}
			}(),
			apiPath:    "/john",
			wantBody:   toJsonFormat(t, john),
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-exist",
			apiPath:    "/john",
			wantBody:   nil, // body確認省略
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := setup(t, tt.db)
			defer c.Close()

			resp, err := c.Do("GET", tt.apiPath, tt.body, "")
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
