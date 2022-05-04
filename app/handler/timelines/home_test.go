package timelines_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/handlertest"
	"yatter-backend-go/app/handler/timelines"
)

// TODO: parameterに関するテスト
func TestHome(t *testing.T) {
	registeredUser := handlertest.AccountData{
		ID:       1,
		UserName: "benben",
	}
	tests := []struct {
		name       string
		db         *handlertest.DBMock
		param      map[string]string
		authUser   string
		wantStatus int
		toTestBody bool
		wantBody   []byte
	}{
		{
			name: "success",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[registeredUser.ID] = registeredUser
				return &handlertest.DBMock{Account: a}
			}(),
			param:      nil,
			authUser:   registeredUser.UserName,
			wantStatus: http.StatusOK,
			toTestBody: true,
			wantBody:   handlertest.ToJsonFormat(t, []object.Status{}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := handlertest.Setup(t, tt.db, timelines.NewRouter)
			defer c.Close()

			resp, err := c.GetWithParamAuth("/home", handlertest.ParamAsURI(tt.param), tt.authUser)
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
