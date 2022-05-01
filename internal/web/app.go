package web

import (
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/processor"
	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/eveisesi/skillz/public"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	csrf "github.com/gobuffalo/mw-csrf"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type Service struct {
	env        skillz.Environment
	baseDomain string
	app        *buffalo.App
	auth       auth.API
	user       user.API
	processor  *processor.Service
	logger     *logrus.Logger
	renderer   *render.Engine
	newrelic   *newrelic.Application
}

const keyAuthenticatedUser = "authenticatedUser"
const keyAuthenticatedUserID = "authenticatedUserID"

const titleSuffix = "|| Eve Is ESI || A Third Party Eve Online App"

func NewService(
	env skillz.Environment,
	sessionName string,
	logger *logrus.Logger,

	auth auth.API,
	user user.API,
	processor *processor.Service,

	renderer *render.Engine,
	newrelic *newrelic.Application,
) *Service {

	var baseDomain = "http://localhost:54400"
	if env == skillz.Production {
		baseDomain = "https://skillboard.eveisesi.space"
	}

	s := &Service{
		env:        env,
		baseDomain: baseDomain,
		auth:       auth,
		user:       user,
		processor:  processor,
		renderer:   renderer,
		logger:     logger,
		newrelic:   newrelic,
	}

	s.app = buffalo.New(buffalo.Options{
		Env:         env.String(),
		SessionName: sessionName,
		WorkerOff:   true,
		Addr:        "0.0.0.0:54400",
	})

	mux := s.app.Muxer()
	mux.Use(s.monitoring)

	s.app.Use(s.setBaseDomain)
	s.app.Use(s.setCurrentUser)
	s.app.GET("/", s.indexHandler)
	s.app.GET("/login", csrf.New(s.loginGetHandler))
	s.app.POST("/login", csrf.New(s.loginPostHandler))
	s.app.GET("/logout", s.logoutHandler)
	s.app.GET("/users/settings", csrf.New(s.authorize(s.userSettingsHandler)))
	s.app.POST("/users/settings", csrf.New(s.authorize(s.postUserSettingsHandler)))
	s.app.DELETE("/users/settings", csrf.New(s.authorize(s.deleteUserSettingsHandler)))
	s.app.GET("/users/{userID}", s.userHandler)

	s.app.ServeFiles("/", http.FS(public.FS())) // serve files from the public directory

	return s

}

func (s *Service) Start() error {
	return s.app.Serve()
}

func (s *Service) flashDanger(c buffalo.Context, msg string) {
	c.Flash().Add("danger", msg)
}

func (s *Service) flashSuccess(c buffalo.Context, msg string) {
	c.Flash().Add("success", msg)
}
