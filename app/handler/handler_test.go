package handler_test

import (
	"net/http"
	"testing"

	"yatter-backend-go/app/handler"
	"yatter-backend-go/app/handler/handlertest"
)

// GETのみ試して, 各ルーティングができているかのみ確認する
func TestHandler(t *testing.T) {
	john := handlertest.AccountData{
		ID:       1,
		UserName: "john",
	}
	johnsStatus := handlertest.StatusData{
		ID:       1,
		Content:  "john's status",
		UserName: john.UserName,
	}
	tests := []struct {
		name       string
		db         *handlertest.DBMock
		apiPath    string
		wantStatus int
	}{
		{
			name: "account",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				return &handlertest.DBMock{Account: a}
			}(),
			apiPath:    "/v1/accounts/john",
			wantStatus: http.StatusOK,
		},
		{
			name: "status",
			db: func() *handlertest.DBMock {
				s := make(handlertest.StatusTableMock)
				s[johnsStatus.ID] = johnsStatus
				return &handlertest.DBMock{Status: s}
			}(),
			apiPath:    "/v1/statuses/1",
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := handlertest.Setup(t, tt.db, handler.NewRouter)
			defer c.Close()

			resp, err := c.Get(tt.apiPath)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// check status code
			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("code expected: %v, returned: %v", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}
