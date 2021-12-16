package server

import "net/http"

// func (s *server) GetOpenIDConfiguration(w http.ResponseWriter, r *http.Request) {
// 	_ = json.NewEncoder(w).Encode(map[string]interface{}{
// 		"jwks_uri": "https://localhost:54405/auth/jwks",
// 	})
// }

// func (s *server) GetJsonWebKeySet(w http.ResponseWriter, r *http.Request) {
// 	_ = json.NewEncoder(w).Encode(s.auth.GetPublicUserJWKS())

// }

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

		s.writeResponse(ctx, w, http.StatusOK, user)
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
