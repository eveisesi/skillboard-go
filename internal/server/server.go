package server

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/graphql"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type server struct {
	port   uint
	env    skillz.Environment
	logger *logrus.Logger

	cache   cache.PageAPI
	graphql graphql.API

	server *http.Server

	templates *template.Template
}

func New(
	port uint,
	env skillz.Environment,
	logger *logrus.Logger,
	cache cache.PageAPI,
	graphql graphql.API,
) *server {

	s := &server{
		port:    port,
		env:     env,
		logger:  logger,
		cache:   cache,
		graphql: graphql,
	}

	s.parseTemplateFiles()

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.buildRouter(),
	}

	return s
}

func (s *server) parseTemplateFiles() {

	var paths []string
	err := filepath.Walk("views", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		paths = append(paths, path)
		return nil
	})
	if err != nil {
		s.logger.WithError(err).Fatal("failed to view dir for templates")
	}

	tmpls, err := template.ParseFiles(paths...)
	if err != nil {
		s.logger.WithError(err).Fatal("failed to parse template files")
	}

	s.templates = tmpls

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
	)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// r.Group(func(r chi.Router) {
	// 	r.Use(
	// 		middleware.SetHeader(headers.ContentType, "application/json"),
	// 	)
	// 	r.Get("/.well-known/openid-configuration", s.GetOpenIDConfiguration)
	// 	r.Get("/auth/jwks", s.GetJsonWebKeySet)
	// })

	r.Get("/", s.handleRenderHomepage)
	r.Get("/characters/{userUUID}", s.handleRenderCharacterPage)

	return r
}
