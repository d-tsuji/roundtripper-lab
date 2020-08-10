package roundtripper_lab

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type retryableRoundTripper struct {
	tr       http.RoundTripper
	attempts int
}

func (t *retryableRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	for attempt := 0; attempt < t.attempts; attempt++ {
		log.Printf("retryableRoundTripper retry: %d\n", attempt+1)
		resp, err = t.tr.RoundTrip(req)

		if err != nil || resp.StatusCode != http.StatusTooManyRequests {
			return resp, err
		}
		// TODO(d-tsuji): We need to take into account the case where the response header contains "Retry-After".

		// Read and close the response body to reuse TCP connections.
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}

	return resp, err
}
