package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/graphql"
	"github.com/eveisesi/skillz/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-http-utils/headers"
	"github.com/sirupsen/logrus"
)

type server struct {
	port   uint
	env    skillz.Environment
	logger *logrus.Logger

	auth    auth.API
	graphql graphql.API
	user    user.API

	server *http.Server
}

func New(
	port uint,
	env skillz.Environment,
	logger *logrus.Logger,
	auth auth.API,
	graphql graphql.API,
	user user.API,
) *server {

	s := &server{
		port:   port,
		env:    env,
		logger: logger,

		auth:    auth,
		graphql: graphql,
		user:    user,
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.buildRouter(),
	}

	return s
}

func (s *server) Run() error {
	entry := s.logger.WithField("service", "server")
	entry.Infof("Starting on Port %d", s.port)

	return s.server.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *server) buildRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Group(func(r chi.Router) {
		r.Use(
			s.requestLogger(s.logger),
			s.cors,
			middleware.SetHeader(headers.ContentType, "application/json"),
			s.authorization,
		)

		// ##### GraphQL Handler #####
		handler := handler.New(s.graphql.ExecutableSchema())

		handler.AddTransport(transport.POST{})
		handler.Use(extension.Introspection{})
		r.Handle("/graphql", handler)
		r.Get("/playground", playground.Handler("GraphQL playground", "/graphql"))

	})

	return r
}
