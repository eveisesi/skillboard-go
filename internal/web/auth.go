package web

import (
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz/internal/errors"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

func (s *Service) logoutHandler(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "You've been logged out successfully")
	return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
}

func (s *Service) loginHandler(c buffalo.Context) error {
	c.Set(defaultTitle())

	var ctx = c.Request().Context()

	code, state := c.Param("code"), c.Param("state")

	if code != "" && state != "" {

		user, err := s.user.Login(ctx, code, state)
		if err != nil {
			s.logger.WithError(err).Error("failed to query user from attempt")
			return errors.NewBuffaloHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to query user from attempt"))
		}

		c.Session().Set(keyAuthenticatedUserID, user.ID)

		return c.Redirect(http.StatusTemporaryRedirect, "userPath()", render.Data{"userID": user.ID.String()})

	}

	attempt, err := s.auth.InitializeAttempt(ctx)
	if err != nil {
		return err
	}

	authorizationURL := s.auth.AuthorizationURI(ctx, attempt.State)

	c.Set("authorizationURL", authorizationURL)

	return c.Render(http.StatusOK, s.renderer.HTML("login/index.plush.html"))
}
