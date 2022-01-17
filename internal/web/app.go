package web

import (
	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	csrf "github.com/gobuffalo/mw-csrf"
	"github.com/sirupsen/logrus"
)

type Service struct {
	app  *buffalo.App
	auth auth.API
	user user.API
	// cache    *cache.PageAPI
	logger   *logrus.Logger
	renderer *render.Engine
}

func NewService(env skillz.Environment,
	sessionName string,
	logger *logrus.Logger,

	auth auth.API,
	user user.API,

	renderer *render.Engine,
) *Service {

	a := buffalo.New(buffalo.Options{
		Env:         env.String(),
		SessionName: sessionName,
		WorkerOff:   true,
		Addr:        "127.0.0.1:54400",
	})

	a.Use(
		csrf.New,
	)

	s := &Service{
		app:      a,
		auth:     auth,
		user:     user,
		renderer: renderer,
		logger:   logger,
	}

	a.GET("/", s.indexHandler)
	a.GET("/login", s.loginHandler)
	a.GET("/robots.txt/", s.robotsHandler)

	return s

}

func (s *Service) Start() error {
	return s.app.Serve()
}
