package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type Server struct {
	logger   *logrus.Logger
	newrelic *newrelic.Application

	auth  auth.API
	users user.API

	http *http.Server
}

type Response struct {
	Success bool
	Error   error
	User    *skillz.User
	Users   []*skillz.User
}

func New(logger *logrus.Logger, newrelic *newrelic.Application, auth auth.API, user user.API) *Server {
	s := &Server{
		logger:   logger,
		newrelic: newrelic,
		auth:     auth,
		users:    user,
	}

	s.http = &http.Server{
		Addr:         ":54400",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		Handler:      s.buildRouter(),
	}

	return s
}

func (s *Server) Start() error {
	s.logger.Infof("Starting Server on Port: %s", s.http.Addr)
	return s.http.ListenAndServe()
}

func (s *Server) buildRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		s.monitoring,
		s.requestLogger(),
		middleware.SetHeader("Content-Type", "application/json"),
	)
	r.Get("/recent", s.handleGetRecent)
	r.Get("/users/{userID}", s.handleGetRecent)
	return r
}

// GracefullyShutdown gracefully shuts down the HTTP API.
func (s *Server) GracefullyShutdown(ctx context.Context) error {
	s.logger.Info("attempting to shutdown server gracefully")
	return s.http.Shutdown(ctx)
}

func (s *Server) writeResponse(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {

	if code != http.StatusOK {
		w.WriteHeader(code)
	}

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}

}

func (s *Server) writeError(ctx context.Context, w http.ResponseWriter, code int, err error) {

	// If err is not nil, actually pass in a map so that the output to the wire is {"error": "text...."} else just let it fall through
	if err != nil {
		s.writeResponse(ctx, w, code, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	s.writeResponse(ctx, w, code, nil)

}
