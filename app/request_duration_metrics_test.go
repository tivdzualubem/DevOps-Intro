package main

import (
	"net/http"
	"strings"
	"testing"
)

func TestMetrics_ExposesRequestDurationHistogram(t *testing.T) {
	srv := newTestServer(t)

	_ = do(t, srv, http.MethodGet, "/health", nil)
	recorder := do(t, srv, http.MethodGet, "/metrics", nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("metrics status: %d", recorder.Code)
	}

	body := recorder.Body.String()

	for _, expected := range []string{
		"# TYPE quicknotes_http_request_duration_seconds histogram",
		`quicknotes_http_request_duration_seconds_bucket{le="+Inf"} 1`,
		"quicknotes_http_request_duration_seconds_count 1",
	} {
		if !strings.Contains(body, expected) {
			t.Errorf("metrics missing %q", expected)
		}
	}
}
