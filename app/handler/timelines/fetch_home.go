package timelines

import (
	"net/http"
)

// Handle request for `GET /v1/timelines/home`
func (h *handler) FetchHome(_ http.ResponseWriter, _ *http.Request) {
}
