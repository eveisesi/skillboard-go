package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-http-utils/headers"
)

func (s *server) handleGetAuthCallback(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	code, state, err := s.parseCodeAndStateFromURL(r.URL)
	if err != nil {
		s.writeError(ctx, w, http.StatusBadRequest, err)
		return
	}

	err = s.user.Login(ctx, code, state)
	if err != nil {
		s.logger.WithError(err).Errorln()
		s.writeError(ctx, w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set(headers.ContentType, "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`
		<!DOCTYPE html>
		<html>
			<title>Athena EVE SSO Auth Callback</title>
			<style>
			body {
				background-color: #585858;
			}
			</style>
			<body>
				<h2>Athena EVE SSO Auth Callback</h2>
			</body>
			<script>
				setTimeout(function() {
					window.close()
				}, 1000)
			</script>
		</html>
	`))

}

func (s *server) parseCodeAndStateFromURL(uri *url.URL) (code, state string, err error) {

	code = uri.Query().Get("code")
	state = uri.Query().Get("state")
	if code == "" || state == "" {
		return "", "", fmt.Errorf("required paramter missing from request")
	}

	return code, state, nil

}
