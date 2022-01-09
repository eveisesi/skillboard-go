package server

import (
	"net/http"

	"github.com/eveisesi/skillz/internal"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	CookieID = "skillz-authed-user-id"
)

// Cors middleware to allow frontend consumption
func (s *server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "600")

		// Just return for OPTIONS
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) authorization(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx = r.Context()

		cookie, err := r.Cookie(CookieID)
		if err == nil {
			user, err := s.user.UserByCookie(ctx, cookie)
			if err != nil {
				s.logger.WithError(err).Error("failed to fetch user by cookie value")
				cookie, err := s.auth.LogoutCookie(ctx)
				if err != nil {
					s.logger.WithError(err).Error("failed to build logout cookie")
					s.writeError(ctx, w, http.StatusInternalServerError, errors.New("failed to process request, please try again later"))
					return
				}

				r.AddCookie(cookie)

				s.writeResponse(ctx, w, http.StatusBadRequest, errors.Wrap(err, "failed to verify user cookie"))
				return
			}

			ctx = internal.ContextWithUser(ctx, user)

			next.ServeHTTP(w, r.WithContext(ctx))

		}
	})

}

// NewStructuredLogger is a constructor for creating a request logger middleware
func (s *server) requestLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&structuredLogger{logger})
}
