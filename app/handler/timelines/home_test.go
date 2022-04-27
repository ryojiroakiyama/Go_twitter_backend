package timelines_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/handlertest"
	"yatter-backend-go/app/handler/timelines"
)

// このテストのみhandlertest内のmockではなく,
// 実際に挙動をベタ書きしたmock(home_mock_tes.go内)でやってみた
func TestHome(t *testing.T) {
	tests := []struct {
		name       string
		prameter   params
		authUser   string
		wantStatus int
		toTestBody bool
		wantBody   []byte
	}{
		{
			name:       "success",
			prameter:   params{},
			authUser:   registeredUser,
			wantStatus: http.StatusOK,
			toTestBody: true,
			wantBody:   handlertest.ToJsonFormat(t, []object.Status{}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := &app.App{Dao: &daoMock{}}
			c := &handlertest.C{
				App:    app,
				Server: httptest.NewServer(timelines.NewRouter(app)),
			}
			defer c.Close()

			resp, err := c.GetWithAuth("/home"+tt.prameter.asURI(), tt.authUser)
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
