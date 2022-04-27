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

const (
	limit = "limit"
)

func TestFollowing(t *testing.T) {
	john := handlertest.AccountData{
		ID:       1,
		UserName: "john",
	}
	benben := handlertest.AccountData{
		ID:       2,
		UserName: "benben",
	}
	sonson := handlertest.AccountData{
		ID:       3,
		UserName: "sonson",
	}
	john_follow_benben := handlertest.RelationShipData{
		UserID:   john.ID,
		TargetID: benben.ID,
	}
	john_follow_sonson := handlertest.RelationShipData{
		UserID:   john.ID,
		TargetID: sonson.ID,
	}
	tests := []struct {
		name       string
		db         *handlertest.DBMock
		param      map[string]string
		username   string
		wantStatus int
		toTestBody bool
		wantBody   []byte
	}{
		{
			name: "success default",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				a[benben.ID] = benben
				a[sonson.ID] = sonson
				r := make(handlertest.RelationShipTableMock)
				r[0] = john_follow_benben
				r[1] = john_follow_sonson
				return &handlertest.DBMock{Account: a, RelationShip: r}
			}(),
			param:      nil,
			username:   john.UserName,
			wantStatus: http.StatusOK,
			toTestBody: true,
			wantBody: handlertest.ToJsonFormat(t,
				[]object.Account{
					{
						ID:       benben.ID,
						Username: benben.UserName,
					},
					{
						ID:       sonson.ID,
						Username: sonson.UserName,
					},
				}),
		},
		{
			name: "success limit=1",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				a[benben.ID] = benben
				a[sonson.ID] = sonson
				r := make(handlertest.RelationShipTableMock)
				r[0] = john_follow_benben
				r[1] = john_follow_sonson
				return &handlertest.DBMock{Account: a, RelationShip: r}
			}(),
			param:      map[string]string{limit: "1"},
			username:   john.UserName,
			wantStatus: http.StatusOK,
			toTestBody: true,
			wantBody: handlertest.ToJsonFormat(t,
				[]object.Account{
					{
						ID:       benben.ID,
						Username: benben.UserName,
					},
				}),
		},
		{
			name: "not_found",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				a[benben.ID] = benben
				r := make(handlertest.RelationShipTableMock)
				r[0] = john_follow_benben
				return &handlertest.DBMock{Account: a, RelationShip: r}
			}(),
			param:      nil,
			username:   benben.UserName,
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := handlertest.Setup(t, tt.db, accounts.NewRouter)
			defer c.Close()

			resp, err := c.GetWithParam("/"+tt.username+"/following", handlertest.ParamAsURI(tt.param))
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
