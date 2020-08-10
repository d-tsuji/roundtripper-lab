package roundtripper_lab

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type loggingRoundTripper struct {
	Transport http.RoundTripper
	Logger    func(string, ...interface{})
}

func (lrt *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if lrt.Logger == nil {
		lrt.Logger = log.Printf
	}

	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		lrt.Logger("Could not dump request: %s", err)
	}
	lrt.Logger("Send request: %s", string(dump))

	res, err := lrt.Transport.RoundTrip(req)

	dump, err = httputil.DumpResponse(res, true)
	if err != nil {
		lrt.Logger("Could not dump response: %s", err)
	}
	lrt.Logger("Received response: %s", string(dump))

	return res, err
}
