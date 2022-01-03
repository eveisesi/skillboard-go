package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/alliance"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/clone"
	"github.com/eveisesi/skillz/internal/contact"
	"github.com/eveisesi/skillz/internal/corporation"
	"github.com/eveisesi/skillz/internal/skill"
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

	auth        auth.API
	alliance    alliance.API
	character   character.API
	clones      clone.API
	contacts    contact.API
	corporation corporation.API
	skills      skill.API
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
	clones clone.API,
	contacts contact.API,
	corporation corporation.API,
	skills skill.API,
	user user.API,
) *server {

	s := &server{
		port:   port,
		env:    env,
		logger: logger,

		auth:        auth,
		alliance:    alliance,
		character:   character,
		clones:      clones,
		contacts:    contacts,
		corporation: corporation,
		skills:      skills,
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

	r.Get("/recent", s.handleGetNewUsersBySP)
	r.Get("/search", s.handleGetUserSearch)

	r.Group(func(r chi.Router) {
		r.Use(
			s.authorization,
		)
		r.Get("/auth", s.handleGetAuth)

		r.Get("/user/{userID}", s.handleGetUserByID)
		r.Get("/user/{userID}/character", s.handleGetUserCharacterByID)
		r.Get("/user/{userID}/refresh", s.handleGetUserByIDRefresh)
		r.Get("/user/{userID}/clones", s.handleGetUserClonesByID)
		r.Get("/user/{userID}/implants", s.handleGetUserImplantsByID)
		r.Get("/user/{userID}/skills/meta", s.handleGetUserSkillMetaByID)
		r.Get("/user/{userID}/skills", s.handleGetUserSkillsByID)
		r.Get("/user/{userID}/queue", s.handleGetUserQueueByID)
		r.Get("/user/{userID}/attributes", s.handleGetUserAttributesByID)
		r.Get("/user/{userID}/flyable", s.handleGetUserFlyableByID)
		r.Get("/user/{userID}/contacts", s.handleGetUserContactsByID)

	})

	return r
}

func (s *server) writeResponse(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {

	if code != http.StatusOK {
		w.WriteHeader(code)
	}

	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.logger.WithError(err).Error("failed to encode data has JSON")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message": "InternalServerError"}`))
		}
	}
}

func (s *server) writeError(ctx context.Context, w http.ResponseWriter, code int, err error) {

	// If err is not nil, actually pass in a map so that the output to the wire is {"error": "text...."} else just let it fall through
	if err != nil {
		s.writeResponse(ctx, w, code, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	s.writeResponse(ctx, w, code, nil)

}
