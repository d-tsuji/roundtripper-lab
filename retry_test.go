package sample

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestRoundTripperRetry(t *testing.T) {
	var counter int64
	testserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t.Log("got request")
		atomic.AddInt64(&counter, 1)
		// Return a retryable error for testing.
		w.WriteHeader(http.StatusTooManyRequests)
	}))

	// Use a retryable client. This will retry up to five times.
	client := &http.Client{
		Transport: &retryableRoundTripper{
			tr:       http.DefaultTransport,
			attempts: 5,
		},
	}

	req, err := http.NewRequest(http.MethodGet, testserver.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("response is wrong, got=%v, want=%v", resp.StatusCode, http.StatusTooManyRequests)
	}

	if counter != 5 {
		t.Errorf("counter is wrong, got=%v, want=%v", counter, 5)
	}
}
