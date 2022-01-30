package web

import (
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/errors"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

func (s *Service) logoutHandler(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "You've been logged out successfully")
	return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
}

func (s *Service) loginGetHandler(c buffalo.Context) error {
	c.Set("title", fmt.Sprintf("Welcome to Skillboard.Evie %s", titleSuffix))

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

	return c.Render(http.StatusOK, s.renderer.HTML("login/index.plush.html"))
}

func (s *Service) loginPostHandler(c buffalo.Context) error {
	var r = c.Request()
	var ctx = r.Context()
	var form = r.Form

	scopes := []string{skillz.ReadSkillsV1.String(), skillz.ReadSkillQueueV1.String()}

	if form.Has("allow_implants") {
		scopes = append(scopes, skillz.ReadImplantsV1.String())
	}

	attempt, err := s.auth.InitializeAttempt(ctx)
	if err != nil {
		return err
	}

	authorizationURL := s.auth.AuthorizationURI(ctx, attempt.State, scopes)

	return c.Redirect(http.StatusFound, authorizationURL)

}
