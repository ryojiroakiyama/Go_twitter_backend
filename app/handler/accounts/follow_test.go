package accounts_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"yatter-backend-go/app/domain/object"
)

func TestFollow(t *testing.T) {
	john := &object.Account{
		Username: "john",
	}
	tests := []struct {
		name       string
		db         *dbMock
		authUser   string
		followUser string
		payload    string
		wantStatus int
		toTestBody bool
		wantBody   []byte
	}{
		{
			name: "try non-exist user",
			db: func() *dbMock {
				a := make(accountTableMock)
				a[john.Username] = *john
				return &dbMock{account: a}
			}(),
			authUser:   "john",
			followUser: "no_such_account",
			payload:    "",
			wantStatus: http.StatusNotFound,
			toTestBody: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := setup(t, tt.db)
			defer c.Close()

			resp, err := c.PostWithAuth("/"+tt.followUser+"/follow", tt.payload, tt.authUser)
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
