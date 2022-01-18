package web

import (
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/eveisesi/skillz/public"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
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
	a.GET("/user/{userID}", s.userHandler)
	// a.GET("/user/{userID}/settings", s.userHandler)

	a.GET("/ping", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, s.renderer.String("pong"))
	})

	a.ServeFiles("/", http.FS(public.FS())) // serve files from the public directory

	return s

}

func (s *Service) Start() error {
	return s.app.Serve()
}
