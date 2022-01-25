package server

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
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

// func (s *server) authorization(next http.Handler) http.Handler {

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		var ctx = r.Context()

// 		authHeader := r.Header.Get("authorization")
// 		if authHeader != "" {
// 			if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
// 				err := errors.New("invalid token, missing token type designator")
// 				LogEntrySetField(ctx, "error", err)
// 				s.writeError(ctx, w, http.StatusUnauthorized, err)
// 				return
// 			}

// 			tokenStr := authHeader[len("bearer "):]
// 			if tokenStr == "null" {
// 				err := errors.New("received null for token, please remove header if token is not available")
// 				LogEntrySetField(ctx, "error", err)
// 				s.writeError(ctx, w, http.StatusUnauthorized, err)
// 				return
// 			}

// 			userID, err := s.auth.UserIDFromToken(ctx, tokenStr)
// 			if err != nil {
// 				LogEntrySetField(ctx, "error", errors.Wrap(err, "failed to parse token for valid userID"))
// 				s.writeError(ctx, w, http.StatusUnauthorized, errors.New("failed to parse token for valid userID"))
// 				return
// 			}

// 			user, err := s.user.User(ctx, userID)
// 			if err != nil {
// 				LogEntrySetField(ctx, "error", errors.Wrap(err, "unknown user for provided token"))
// 				s.writeError(ctx, w, http.StatusUnauthorized, errors.New("unknown user for provided token"))
// 				return
// 			}

// 			user.Settings, err = s.user.UserSettings(ctx, user.ID)
// 			if err != nil {
// 				LogEntrySetField(ctx, "error", errors.Wrap(err, "unknown user for provided token"))
// 			}

// 			ctx = internal.ContextWithUser(ctx, user)
// 		}

// 		next.ServeHTTP(w, r.WithContext(ctx))

// 	})

// }

// func (s *server) hasPermission(permission user.Permission, handler http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		var ctx = r.Context()

// 		targetIDStr := chi.URLParam(r, "userID")
// 		if targetIDStr == "" {
// 			handler.ServeHTTP(w, r)
// 			return
// 		}

// 		targetID, err := uuid.FromString(targetIDStr)
// 		if err != nil {
// 			handler.ServeHTTP(w, r)
// 			return
// 		}

// 		tokenUser := internal.UserFromContext(ctx)
// 		if tokenUser != nil && tokenUser.ID == targetID {
// 			handler.ServeHTTP(w, r)
// 			return
// 		}

// 		target, err := s.user.User(ctx, targetID)
// 		if err != nil {
// 			handler.ServeHTTP(w, r)
// 			return
// 		}

// 		if target.Settings == nil {
// 			settings, err := s.user.UserSettings(ctx, targetID)
// 			if err != nil {
// 				handler.ServeHTTP(w, r)
// 				return
// 			}

// 			if settings == nil {
// 				handler.ServeHTTP(w, r)
// 				return
// 			}

// 			target.Settings = settings

// 		}

// 		settings := target.Settings

// 		switch p := permission; {
// 		case p == user.PermissionHideClones && settings.HideClones:
// 			fallthrough
// 		case p == user.PermissionHideQueue && settings.HideQueue:
// 			fallthrough
// 		case p == user.PermissionHideShips && settings.HideShips:
// 			fallthrough
// 		case p == user.PermissionHideStandings && settings.HideStandings:
// 			s.writeResponse(ctx, w, http.StatusForbidden, errors.New(http.StatusText(http.StatusForbidden)))
// 			return
// 		}

// 		handler.ServeHTTP(w, r)

// 	}
// }

// NewStructuredLogger is a constructor for creating a request logger middleware
func (s *server) requestLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&structuredLogger{logger})
}
