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

func TestPublic(t *testing.T) {
	johnsStatus := handlertest.StatusData{
		ID:       1,
		Content:  "john's status",
		UserName: "john",
	}
	benbensStatus := handlertest.StatusData{
		ID:       2,
		Content:  "benben's status",
		UserName: "benben",
	}
	sonsonsStatus := handlertest.StatusData{
		ID:       3,
		Content:  "sonson's status",
		UserName: "sonson",
	}
	tests := []struct {
		name       string
		db         *handlertest.DBMock
		prameter   params
		wantStatus int
		toTestBody bool
		wantBody   []byte
	}{
		{
			name: "success",
			db: func() *handlertest.DBMock {
				s := make(handlertest.StatusTableMock)
				s[johnsStatus.ID] = johnsStatus
				s[benbensStatus.ID] = benbensStatus
				s[sonsonsStatus.ID] = sonsonsStatus
				return &handlertest.DBMock{Status: s}
			}(),
			prameter:   params{},
			wantStatus: http.StatusOK,
			toTestBody: true,
			wantBody: handlertest.ToJsonFormat(t, []object.Status{
				{
					ID:      johnsStatus.ID,
					Content: johnsStatus.Content,
					Account: &object.Account{
						Username: johnsStatus.UserName,
					},
				},
				{
					ID:      benbensStatus.ID,
					Content: benbensStatus.Content,
					Account: &object.Account{
						Username: benbensStatus.UserName,
					},
				},
				{
					ID:      sonsonsStatus.ID,
					Content: sonsonsStatus.Content,
					Account: &object.Account{
						Username: sonsonsStatus.UserName,
					},
				},
			}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := handlertest.Setup(t, tt.db, timelines.NewRouter)
			defer c.Close()

			resp, err := c.Get("/public" + tt.prameter.asURI())
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
