package statuses_test

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"testing"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/handlertest"
	"yatter-backend-go/app/handler/statuses"
)

func TestFetch(t *testing.T) {
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
		wantStatus int
		toTestBody bool
		wantBody   []byte
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
			wantStatus: http.StatusOK,
			toTestBody: true,
			wantBody: handlertest.ToJsonFormat(t, object.Status{
				ID:      johnsStatus.ID,
				Content: johnsStatus.Content,
				Account: &object.Account{
					Username: johnsStatus.UserName,
				},
			}),
		},
		{
			name:       "non exist",
			db:         nil,
			id:         johnsStatus.ID,
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := handlertest.Setup(t, tt.db, statuses.NewRouter)
			defer c.Close()

			resp, err := c.Get("/" + strconv.FormatInt(tt.id, 10))
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
