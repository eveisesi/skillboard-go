package web

import (
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/gertd/go-pluralize"
	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

var userPageTitle = func(name string) (string, string) {
	name = pluralize.NewClient().Plural(name)
	return "title", fmt.Sprintf("%s Skillboard %s", name, titleSuffix)
}

var userSettingsPageTitle = func(name string) (string, string) {
	name = pluralize.NewClient().Plural(name)
	return "title", fmt.Sprintf("%s Settings %s", name, titleSuffix)
}

func (s *Service) userHandler(c buffalo.Context) error {
	var ctx = c.Request().Context()

	userID, err := uuid.FromString(c.Param("userID"))
	if err != nil {
		c.Flash().Add("danger", "invalid id supplied for userID")
		return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
	}

	u, err := s.user.User(ctx, userID)
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return c.Error(http.StatusInternalServerError, err)
	}

	if errors.Is(err, user.ErrUserNotFound) {
		c.Flash().Add("danger", err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
	}

	settings := u.Settings
	if settings != nil {
		if settings.Visibility == skillz.VisibilityPrivate {
			sessionUser := c.Data()[keyAuthenticatedUser]
			if sessionUser == nil {
				s.flashDanger(c, "User Not Found")
				return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
			}

			authenticatedUser := sessionUser.(*skillz.User)
			if authenticatedUser.ID != userID {
				s.flashDanger(c, "User Not Found")
				return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
			}
		}

		if settings.Visibility == skillz.VisibilityToken {

			token := c.Param("token")
			if token != "" {

				if token != "" && token != settings.VisibilityToken {
					s.flashDanger(c, "User Not Found")
					return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
				}
			} else {
				sessionUser := c.Data()[keyAuthenticatedUser]
				if sessionUser == nil {
					s.flashDanger(c, "User Not Found")
					return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
				}

				authenticatedUser, ok := sessionUser.(*skillz.User)
				if !ok {
					s.flashDanger(c, "User Not Found")
					return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
				}

				if authenticatedUser.ID != userID {
					s.flashDanger(c, "User Not Found")
					return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
				}
			}

		}
	}

	c.Set(defaultTitle())
	if u.IsNew {
		return c.Render(http.StatusOK, s.renderer.HTML("user/welcome.plush.html"))
	}

	u, err = s.user.LoadUserAll(ctx, userID)
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return c.Error(http.StatusInternalServerError, err)
	}
	if u.Character != nil {
		c.Set(userPageTitle(u.Character.Name))
	}

	c.Set("user", u)

	return c.Render(http.StatusOK, s.renderer.HTML("user/index.plush.html"))
}

func (s *Service) userSettingsHandler(c buffalo.Context) error {

	user := c.Data()[keyAuthenticatedUser].(*skillz.User)
	if user == nil {
		s.logger.Debug("user is missing from session")
		return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
	}

	c.Set(userSettingsPageTitle(user.Character.Name))
	c.Set("checked", func(b bool) string {
		o := ""
		if b {
			o = `checked="on"`
		}
		return o
	})
	c.Set("visibilities", skillz.AllVisibilities)
	return c.Render(http.StatusOK, s.renderer.HTML("user/settings.plush.html"))
}

func (s *Service) postUserSettingsHandler(c buffalo.Context) error {

	var ctx = c.Request().Context()

	user := c.Data()[keyAuthenticatedUser].(*skillz.User)
	if user == nil {
		return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
	}

	settings := new(skillz.UserSettings)
	settings.VisibilityToken = user.Settings.VisibilityToken

	err := c.Bind(settings)
	if err != nil {
		s.flashDanger(c, "failed to process form. Please try again")
		return s.userSettingsHandler(c)
	}

	if !settings.Visibility.Valid() {
		s.flashDanger(c, "Invalid value for Visibilty. Please try again")
		return s.userSettingsHandler(c)
	}

	err = s.user.CreateUserSettings(ctx, user.ID, settings)
	if err != nil {
		return err
	}

	user.Settings = settings

	routeInfo, _ := s.app.Routes().Lookup("userPath")
	if routeInfo != nil {
		html, _ := routeInfo.BuildPathHelper()(map[string]interface{}{"userID": user.ID})
		if string(html) != "" {
			err = s.cache.BustPage(ctx, string(html))
			if err != nil {
				s.logger.WithError(err).Error("failed to bust page cache")
			}
		}
	}

	s.flashSuccess(c, "Settings updated successfully")
	return c.Redirect(http.StatusFound, "usersSettingsPath()")

}
