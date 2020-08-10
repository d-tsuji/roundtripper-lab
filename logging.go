package sample

import (
	"log"
	"net/http"
	"net/http/httputil"
)

// NewLoggingRoundTripper gets the http.RoundTripper that can log the HTTP request / response.
func NewLoggingRoundTripper(transport http.RoundTripper, logger func(string, ...interface{})) http.RoundTripper {
	return &loggingRoundTripper{transport: transport, logger: logger}
}

type loggingRoundTripper struct {
	transport http.RoundTripper
	logger    func(string, ...interface{})
}

func (lrt *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if lrt.logger == nil {
		lrt.logger = log.Printf
	}

	if dump, err := httputil.DumpRequestOut(req, true); err == nil {
		lrt.logger("Send request: %s", string(dump))
	}

	resp, err := lrt.transport.RoundTrip(req)
	if resp != nil {
		if dump, err := httputil.DumpResponse(resp, true); err == nil {
			lrt.logger("Received response: %s", string(dump))
		}
	}

	return resp, err
}
