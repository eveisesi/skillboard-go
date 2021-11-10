package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/eveisesi/skillz"
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
	// loaders dataloaders.Service

	user user.Service

	server *http.Server
}

func New(
	port uint,
	env skillz.Environment,
	// newrelic *newrelic.Application,
	logger *logrus.Logger,
	// loaders dataloaders.Service,

) *server {

	s := &server{
		port:   port,
		env:    env,
		logger: logger,
		// loaders:     loaders,
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.buildRouter(),
	}

	return s
}

func (s *server) Run() error {
	s.logger.WithField("service", "server").Infof("Starting on Port %d", s.port)
	return s.server.ListenAndServe()
}

// GracefullyShutdown gracefully shuts down the HTTP API.
func (s *server) GracefullyShutdown(ctx context.Context) error {
	s.logger.Info("attempting to shutdown server gracefully")
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

	r.Get("/playground", playground.Handler("GraphQL playground", "/graphql"))

	// r.Group(func(r chi.Router) {
	// 	r.Use(s.authorization)

	// 	// ##### GraphQL Handler #####
	// 	handler := handler.New(
	// 		generated.NewExecutableSchema(
	// 			generated.Config{
	// 				Resolvers: resolvers.New(
	// 					s.logger,
	// 				),
	// 			},
	// 		),
	// 	)
	// 	handler.AddTransport(transport.POST{})
	// 	handler.AddTransport(transport.MultipartForm{})
	// 	handler.Use(extension.Introspection{})
	// 	handler.SetQueryCache(lru.New(1000))
	// 	handler.Use(extension.AutomaticPersistedQuery{
	// 		Cache: lru.New(100),
	// 	})
	// 	r.Handle("/graphql", handler)
	// })

	return r
}
