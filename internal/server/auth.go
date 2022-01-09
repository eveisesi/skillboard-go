package server

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

func (s *server) handleGetAuth(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code != "" && state != "" {
		user, err := s.user.Login(ctx, code, state)
		if err != nil {
			s.writeError(ctx, w, http.StatusBadRequest, err)
			return
		}

		cookie, err := s.auth.UserCookie(ctx, user.ID)
		if err != nil {
			s.writeError(ctx, w, http.StatusInternalServerError, err)
			return
		}

		spew.Dump(cookie)

		r.AddCookie(cookie)

		s.writeResponse(ctx, w, http.StatusOK, user)
		return
	}

	attempt, err := s.auth.InitializeAttempt(ctx)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	spew.Dump(attempt)

	s.writeResponse(ctx, w, http.StatusOK, map[string]interface{}{
		"url": s.auth.AuthorizationURI(ctx, attempt.State),
	})

}

func (s *server) handleGetAuthLogout(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	cookie, err := s.auth.LogoutCookie(ctx)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	r.AddCookie(cookie)

	s.writeResponse(ctx, w, http.StatusNoContent, nil)

}
