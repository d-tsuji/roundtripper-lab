package sample

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func TestRoundTripperLogging(t *testing.T) {
	testserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	}))

	// For example, using rs/zerolog as structured logging would look like this.
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	client := &http.Client{
		Transport: &loggingRoundTripper{
			transport: http.DefaultTransport,
			logger:    log.Printf,
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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("response is wrong, got=%v, want=%v", resp.StatusCode, http.StatusOK)
	}
}
