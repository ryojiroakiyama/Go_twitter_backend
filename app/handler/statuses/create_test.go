package statuses_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/handlertest"
	"yatter-backend-go/app/handler/statuses"
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
		authUser    string
		wantStatus  int
		toTestBody  bool
		wantBody    []byte
	}{
		{
			name: "success",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				return &handlertest.DBMock{Account: a}
			}(),
			postPayload: `{"status":"new status"}`,
			authUser:    john.UserName,
			wantStatus:  http.StatusOK,
			toTestBody:  true,
			wantBody: handlertest.ToJsonFormat(t, object.Status{
				ID:      1, //IDは1から始まるので
				Content: "new status",
				Account: &object.Account{
					Username: john.UserName,
				},
			}),
		},
		{
			name: "bad request",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				return &handlertest.DBMock{Account: a}
			}(),
			postPayload: "notjson",
			authUser:    john.UserName,
			wantStatus:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := handlertest.Setup(t, tt.db, statuses.NewRouter)
			defer c.Close()

			resp, err := c.PostJsonWithAuth("/", tt.postPayload, tt.authUser)
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
