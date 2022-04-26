package accounts_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"yatter-backend-go/app/domain/object"
)

func TestFetch(t *testing.T) {
	john := accountData{
		id:       1,
		username: "benben",
	}
	tests := []struct {
		name       string
		db         *dbMock
		username   string
		wantStatus int
		toTestBody bool
		wantBody   []byte
	}{
		{
			name: "success",
			db: func() *dbMock {
				a := make(accountTableMock)
				a[john.id] = john
				return &dbMock{account: a}
			}(),
			username:   john.username,
			toTestBody: true,
			wantBody:   toJsonFormat(t, object.Account{Username: john.username}),
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-exist",
			username:   "john",
			wantStatus: http.StatusNotFound,
			toTestBody: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := setup(t, tt.db)
			defer c.Close()

			resp, err := c.Get("/" + tt.username)
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
