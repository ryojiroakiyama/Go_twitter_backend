package statuses_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"yatter-backend-go/app/domain/object"
)

func TestCreate(t *testing.T) {
	john := accountData{
		id:       1,
		username: "john",
	}
	//johnsStatus := statusData{
	//	id:       1,
	//	content:  "john's status",
	//	username: john.username,
	//}
	tests := []struct {
		name        string
		db          *dbMock
		postPayload string
		authUser    string
		wantStatus  int
		toTestBody  bool
		wantBody    []byte
	}{
		{
			name: "success",
			db: func() *dbMock {
				a := make(accountTableMock)
				a[john.id] = john
				return &dbMock{account: a}
			}(),
			postPayload: `{"status":"new status"}`,
			authUser:    john.username,
			wantBody: toJsonFormat(t, object.Status{
				ID:      1, //IDは1から始まるので
				Content: "new status",
				Account: &object.Account{
					Username: john.username,
				},
			}),
			wantStatus: http.StatusOK,
		},
		//{
		//	name: "create duplicate status",
		//	db: func() *dbMock {
		//		a := make(accountTableMock)
		//		s := make(statusTableMock)
		//		a[john.Username] = *john
		//		s[1] = *johnStatus1
		//		return &dbMock{account: a, status: s}
		//	}(),
		//	method:       "POST",
		//	apiPath:      "/v1/statuses",
		//	body:         bytes.NewReader([]byte(`{"status":"johnStatus"}`)),
		//	authUserName: john.Username,
		//	wantBody:     toJsonFormat(t, johnStatus2),
		//	wantStatus:   http.StatusOK,
		//},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := setup(t, tt.db)
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
