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

func TestCreate(t *testing.T) {
	john := handlertest.AccountData{
		ID:       1,
		UserName: "john",
	}
	tests := []struct {
		name        string
		db          *handlertest.DBMock
		postPayload string
		wantStatus  int
		toTestBody  bool
		wantBody    []byte
	}{
		{
			name:        "success",
			db:          nil, // 空のDBがセットされる
			postPayload: `{"username":"john"}`,
			wantStatus:  http.StatusOK,
			toTestBody:  true,
			wantBody: handlertest.ToJsonFormat(t,
				object.Account{Username: john.UserName}),
		},
		{
			name: "duplicate",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				return &handlertest.DBMock{Account: a}
			}(),
			postPayload: `{"username":"john"}`,
			wantStatus:  http.StatusConflict,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := handlertest.Setup(t, tt.db, accounts.NewRouter)
			defer c.Close()

			resp, err := c.PostJSON("/", tt.postPayload)
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
