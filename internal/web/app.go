package web

import (
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/eveisesi/skillz/public"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	csrf "github.com/gobuffalo/mw-csrf"
	"github.com/sirupsen/logrus"
)

type Service struct {
	env        skillz.Environment
	baseDomain string
	app        *buffalo.App
	auth       auth.API
	user       user.API
	logger     *logrus.Logger
	renderer   *render.Engine
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

	renderer *render.Engine,
) *Service {

	var baseDomain = "https://skillboard.local"
	if env == skillz.Production {
		baseDomain = "https://skillboard.eveisesi.space"
	}

	s := &Service{
		env:        env,
		baseDomain: baseDomain,
		auth:       auth,
		user:       user,
		renderer:   renderer,
		logger:     logger,
	}

	s.app = buffalo.New(buffalo.Options{
		Env:         env.String(),
		SessionName: sessionName,
		WorkerOff:   true,
		Addr:        "0.0.0.0:54400",
	})

	s.app.Use(s.setBaseDomain)
	s.app.Use(s.setCurrentUser)
	s.app.GET("/", s.indexHandler)
	s.app.GET("/login", s.loginHandler)
	s.app.GET("/logout", s.logoutHandler)
	s.app.GET("/users/settings", csrf.New(s.authorize(s.userSettingsHandler)))
	s.app.PUT("/users/settings", csrf.New(s.authorize(s.postUserSettingsHandler)))
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
