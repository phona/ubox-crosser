package server

import (
	"net/http"
	"testing"
)

// TestStartAdminHTTP_InvalidAddrReturnsError proves the helper surfaces a
// bind error to its caller rather than swallowing it. cmd/server relies on
// this to log the failure through logrus without crashing the proxy loop.
func TestStartAdminHTTP_InvalidAddrReturnsError(t *testing.T) {
	mux := http.NewServeMux()
	// "not-a-valid-addr" has no colon → http.Server rejects at Listen time.
	if err := StartAdminHTTP("not-a-valid-addr", mux); err == nil {
		t.Fatal("want non-nil error for invalid addr, got nil")
	}
}
