package roundtripper

import (
	"net/http"

	"github.com/go-http-utils/headers"
)

func UserAgent(userAgent string, original http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		r.Header.Set(headers.UserAgent, userAgent)
		if original == nil {
			original = http.DefaultTransport
		}

		return original.RoundTrip(r)
	})
}

type roundTripperFunc func(r *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
