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

func TestUnFollow(t *testing.T) {
	john := handlertest.AccountData{
		ID:       1,
		UserName: "john",
	}
	benben := handlertest.AccountData{
		ID:       2,
		UserName: "benben",
	}
	relation := handlertest.RelationShipData{
		UserID:   john.ID,
		TargetID: benben.ID,
	}
	tests := []struct {
		name         string
		db           *handlertest.DBMock
		authUser     string
		unFollowUser string
		wantStatus   int
		toTestBody   bool
		wantBody     []byte
	}{
		{
			name: "success",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				a[benben.ID] = benben
				r := make(handlertest.RelationShipTableMock)
				r[0] = relation
				return &handlertest.DBMock{Account: a, RelationShip: r}
			}(),
			authUser:     john.UserName,
			unFollowUser: benben.UserName,
			wantStatus:   http.StatusOK,
			toTestBody:   true,
			wantBody: handlertest.ToJsonFormat(t, &object.Relationship{
				TargetID:  benben.ID,
				Following: false,
				FllowedBy: false,
			}),
		},
		{
			name: "non-exist account",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				a[benben.ID] = benben
				return &handlertest.DBMock{Account: a}
			}(),
			authUser:     john.UserName,
			unFollowUser: "no such account",
			wantStatus:   http.StatusNotFound,
		},
		{ // これでエラーにはならないという意味でのテスト
			name: "already not following",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				a[benben.ID] = benben
				return &handlertest.DBMock{Account: a}
			}(),
			authUser:     john.UserName,
			unFollowUser: benben.UserName,
			wantStatus:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := handlertest.Setup(t, tt.db, accounts.NewRouter)
			defer c.Close()

			resp, err := c.PostJsonWithAuth("/"+tt.unFollowUser+"/unfollow", "", tt.authUser)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// check status code
			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("code expected: %v, returned: %v", tt.wantStatus, resp.StatusCode)
			}

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			// check response body
			if tt.toTestBody {
				if bytes.Compare(responseBody, tt.wantBody) != 0 {
					t.Fatalf("body \nexpected: [%v], \nreturned: [%v]", string(tt.wantBody), string(responseBody))
				}
			}
		})
	}
}
