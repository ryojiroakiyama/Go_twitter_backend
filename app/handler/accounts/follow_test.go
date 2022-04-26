package accounts_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/accounts"
	"yatter-backend-go/app/handler/handlertest"
)

func TestFollow(t *testing.T) {
	john := handlertest.AccountData{
		ID:       1,
		UserName: "john",
	}
	benben := handlertest.AccountData{
		ID:       2,
		UserName: "benben",
	}
	tests := []struct {
		name       string
		db         *handlertest.DBMock
		authUser   string
		followUser string
		wantStatus int
		toTestBody bool
		wantBody   []byte
	}{
		{
			name: "success",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				a[benben.ID] = benben
				return &handlertest.DBMock{Account: a}
			}(),
			authUser:   john.UserName,
			followUser: benben.UserName,
			wantStatus: http.StatusOK,
			toTestBody: true,
			wantBody: handlertest.ToJsonFormat(t,
				object.Relationship{
					TargetID:  benben.ID,
					Following: true,
					FllowedBy: false}),
		},
		{
			name: "non-exist user to follow",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				return &handlertest.DBMock{Account: a}
			}(),
			authUser:   john.UserName,
			followUser: "no_such_account",
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := handlertest.Setup(t, tt.db, accounts.NewRouter)
			defer c.Close()

			resp, err := c.PostJsonWithAuth("/"+tt.followUser+"/follow", "", tt.authUser)
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
