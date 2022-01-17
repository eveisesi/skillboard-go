package web

import (
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/eveisesi/skillz/internal/errors"
	"github.com/gobuffalo/buffalo"
)

func (s *Service) loginHandler(c buffalo.Context) error {
	c.Set(defaultTitle())

	var ctx = c.Request().Context()

	code, state := c.Param("code"), c.Param("state")

	if code != "" && state != "" {

		user, err := s.user.Login(ctx, state, code)
		if err != nil {
			return errors.NewBuffaloHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to query user from attempt"))
		}

		spew.Dump(user)
		return nil

	}

	attempt, err := s.auth.InitializeAttempt(ctx)
	if err != nil {
		return err
	}

	authorizationURL := s.auth.AuthorizationURI(ctx, attempt.State)

	c.Set("authorizationURL", authorizationURL)

	return c.Render(http.StatusOK, s.renderer.HTML("login/index.plush.html"))
}
