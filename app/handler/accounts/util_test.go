package accounts_test

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

func toJsonFormat(t *testing.T, body interface{}) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}
	return out
}
