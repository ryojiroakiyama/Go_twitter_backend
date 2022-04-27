package statuses_test

import (
	"net/http"
	"strconv"
	"testing"
	"yatter-backend-go/app/handler/handlertest"
	"yatter-backend-go/app/handler/statuses"
)

func TestDelete(t *testing.T) {
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
		id         int64
		authUser   string
		wantStatus int
	}{
		{
			name: "success",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				s := make(handlertest.StatusTableMock)
				s[johnsStatus.ID] = johnsStatus
				return &handlertest.DBMock{Account: a, Status: s}
			}(),
			id:         johnsStatus.ID,
			authUser:   john.UserName,
			wantStatus: http.StatusOK,
		},
		{
			name: "wrong authentication",
			db: func() *handlertest.DBMock {
				a := make(handlertest.AccountTableMock)
				a[john.ID] = john
				s := make(handlertest.StatusTableMock)
				s[johnsStatus.ID] = johnsStatus
				return &handlertest.DBMock{Account: a, Status: s}
			}(),
			id:         johnsStatus.ID,
			authUser:   "no such user",
			wantStatus: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := handlertest.Setup(t, tt.db, statuses.NewRouter)
			defer c.Close()

			resp, err := c.DeleteWithAuth("/"+strconv.FormatInt(tt.id, 10), tt.authUser)
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
