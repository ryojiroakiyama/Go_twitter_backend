package accounts_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"yatter-backend-go/app/domain/object"
)

func TestFollow(t *testing.T) {
	john := accountData{
		id:       1,
		username: "john",
	}
	benben := accountData{
		id:       2,
		username: "benben",
	}
	tests := []struct {
		name       string
		db         *dbMock
		authUser   string
		followUser string
		wantStatus int
		toTestBody bool
		wantBody   []byte
	}{
		{
			name: "success",
			db: func() *dbMock {
				a := make(accountTableMock)
				a[john.id] = john
				a[benben.id] = benben
				return &dbMock{account: a}
			}(),
			authUser:   john.username,
			followUser: benben.username,
			wantStatus: http.StatusOK,
			toTestBody: true,
			wantBody: toJsonFormat(t,
				object.Relationship{
					TargetID:  benben.id,
					Following: true,
					FllowedBy: false}),
		},
		{
			name: "non-exist user to follow",
			db: func() *dbMock {
				a := make(accountTableMock)
				a[john.id] = john
				return &dbMock{account: a}
			}(),
			authUser:   john.username,
			followUser: "no_such_account",
			wantStatus: http.StatusNotFound,
			toTestBody: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := setup(t, tt.db)
			defer c.Close()

			resp, err := c.PostWithAuth("/"+tt.followUser+"/follow", "", tt.authUser)
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
