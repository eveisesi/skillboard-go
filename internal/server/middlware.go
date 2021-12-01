package server

import (
	"context"
	"net/http"

	"github.com/eveisesi/skillz/internal"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gofrs/uuid"
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

type delayedAuthorizationWriter struct {
	http.ResponseWriter
	sessionID uuid.UUID
	cookieFn  func(ctx context.Context, cookie *http.Cookie, userID uuid.UUID) (*http.Cookie, error)
	logger    *logrus.Logger
}

func (d *delayedAuthorizationWriter) Write(data []byte) (int, error) {

	d.logger.Info("checking session")

	userID := internal.CacheGet(d.sessionID)
	if userID == uuid.Nil {
		return d.ResponseWriter.Write(data)
	}

	defer internal.CacheDelete(d.sessionID)

	cookie, err := d.cookieFn(context.Background(), newCookie(), userID)
	if err != nil {
		d.logger.WithError(err).Error("failed to add value to cookie")
		return d.ResponseWriter.Write(data)
	}

	if v := cookie.String(); v != "" {
		d.logger.WithField("cookieLen", len(v)).Info("setting cookie header")
		d.ResponseWriter.Header().Set("Set-Cookie", v)
		return d.ResponseWriter.Write(data)
	} else {
		d.logger.WithField("userID", userID).Error("failed to set cookie due to empty value")
	}

	return d.ResponseWriter.Write(data)

}

func (s *server) authorization(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var ctx = r.Context()

		cookie, err := r.Cookie(CookieID)
		if err == nil {
			user, err := s.user.UserByCookie(ctx, cookie)
			if err != nil {
				s.logger.WithError(err).Error("failed to fetch user by cookie value")
				http.SetCookie(w, &http.Cookie{
					Name:   CookieID,
					MaxAge: -1,
				})

				next.ServeHTTP(w, r)

			} else {
				ctx = internal.ContextWithUser(ctx, user)

				next.ServeHTTP(w, r.WithContext(ctx))
			}

		} else {
			sessionID := uuid.Must(uuid.NewV4())
			internal.CacheSet(sessionID, uuid.Nil)
			ctx = internal.ContextWithSessionID(ctx, sessionID)
			s.logger.WithField("sessionID", sessionID.String()).Info("sessionID generated, cached, and stored on context")

			delayedWriter := &delayedAuthorizationWriter{
				w,
				sessionID,
				s.auth.CookieForUserID,
				s.logger,
			}

			next.ServeHTTP(delayedWriter, r.WithContext(ctx))

		}

	})

}

func newCookie() *http.Cookie {

	return &http.Cookie{
		Name: CookieID,
		// HttpOnly: true,
		// Domain: "skillboard",
		// MaxAge: 50000,
		// Path:   "/",
	}

}

// func (s *server) authorization(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		var ctx = r.Context()

// 		authHeader := r.Header.Get("authorization")

// 		if authHeader != "" {
// 			if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
// 				s.writeError(ctx, w, http.StatusUnauthorized, fmt.Errorf("missing or invalid token"))
// 				return
// 			}

// 			var prefixes = []string{`bearer `, `Bearer `}
// 			for _, prefix := range prefixes {
// 				authHeader = strings.TrimPrefix(authHeader, prefix)
// 			}

// 			token, err := s.auth.ParseAndVerifyUserToken(ctx, authHeader)
// 			if err != nil {
// 				s.writeError(ctx, w, http.StatusUnauthorized, fmt.Errorf("failed to validate token: %w", err))
// 				return
// 			}

// 			user, err := s.user.UserFromToken(ctx, token)
// 			if err != nil {
// 				s.writeError(ctx, w, http.StatusBadRequest, err)
// 				return
// 			}

// 			ctx = internal.ContextWithToken(ctx, token)
// 			ctx = internal.ContextWithUser(ctx, user)
// 			ctx = internal.ContextWithUpdateable(ctx, false)
// 		}

// 		next.ServeHTTP(w, r.WithContext(ctx))

// 	})
// }

// NewStructuredLogger is a constructor for creating a request logger middleware
func (s *server) requestLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&structuredLogger{logger})
}

// func closeRequestBody(ctx context.Context, r *http.Request) {
// 	err := r.Body.Close()
// 	if err != nil {
// 		// newrelic.FromContext(ctx).NoticeError(err)

// 	}
// }

// func (s *server) writeResponse(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {

// 	if code != http.StatusOK {
// 		w.WriteHeader(code)
// 	}

// 	if data != nil {
// 		_ = json.NewEncoder(w).Encode(data)
// 	}
// }

// func (s *server) writeError(ctx context.Context, w http.ResponseWriter, code int, err error) {

// 	// If err is not nil, actually pass in a map so that the output to the wire is {"error": "text...."} else just let it fall through
// 	if err != nil {
// 		LogEntrySetField(ctx, "error", err)
// 		s.writeResponse(ctx, w, code, map[string]interface{}{
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	s.writeResponse(ctx, w, code, nil)

// }
