package server

import (
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
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

		token, err := s.auth.UserToken(ctx, user.ID)
		if err != nil {
			s.writeError(ctx, w, http.StatusInternalServerError, errors.Wrap(err, "failed to generate token"))
			return
		}

		s.writeResponse(ctx, w, http.StatusOK, struct {
			User  *skillz.User `json:"user"`
			Token string       `json:"token"`
		}{
			user, token,
		})

		return
	}

	attempt, err := s.auth.InitializeAttempt(ctx)
	if err != nil {
		s.writeError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	s.writeResponse(ctx, w, http.StatusOK, map[string]interface{}{
		"url": s.auth.AuthorizationURI(ctx, attempt.State),
	})

}
