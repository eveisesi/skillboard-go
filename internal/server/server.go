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
	"github.com/eveisesi/skillz/internal/alliance"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/corporation"
	resolvers "github.com/eveisesi/skillz/internal/server/gql"
	"github.com/eveisesi/skillz/internal/server/gql/dataloaders"
	"github.com/eveisesi/skillz/internal/server/gql/generated"
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

	alliance    alliance.API
	auth        auth.API
	character   character.API
	corporation corporation.API
	dataloaders dataloaders.API
	user        user.API

	server *http.Server
}

func New(

	port uint,
	env skillz.Environment,
	logger *logrus.Logger,

	alliance alliance.API,
	auth auth.API,
	character character.API,
	corporation corporation.API,
	dataloaders dataloaders.API,
	user user.API,
) *server {

	s := &server{
		port:        port,
		env:         env,
		logger:      logger,
		alliance:    alliance,
		auth:        auth,
		character:   character,
		corporation: corporation,
		dataloaders: dataloaders,
		user:        user,
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

	r.Use(
		s.requestLogger(s.logger),
		s.cors,
		middleware.SetHeader(headers.ContentType, "application/json"),
	)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/auth/callback", s.handleGetAuthCallback)

	r.Get("/playground", playground.Handler("GraphQL playground", "/graphql"))

	r.Group(func(r chi.Router) {
		// r.Use(s.authorization)

		// ##### GraphQL Handler #####
		handler := handler.New(
			generated.NewExecutableSchema(
				generated.Config{
					Resolvers: resolvers.New(
						s.alliance,
						s.auth,
						s.character,
						s.corporation,
						s.dataloaders,
						s.user,
					),
				},
			),
		)
		handler.AddTransport(transport.POST{})

		if s.env != skillz.Production {
			handler.Use(extension.Introspection{})
		}

		// handler.SetQueryCache(lru.New(1000))
		// handler.Use(extension.AutomaticPersistedQuery{
		// 	Cache: lru.New(100),
		// })
		r.Handle("/graphql", handler)
	})

	return r
}
