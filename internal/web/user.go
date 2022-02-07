package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/gertd/go-pluralize"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

var printer = message.NewPrinter(language.English)

func (s *Service) userPageMeta(ctx context.Context, user *skillz.User) (string, *Meta) {
	var meta = new(Meta)
	if user.Character != nil {
		name := pluralize.NewClient().Plural(user.Character.Name)
		meta.Title = fmt.Sprintf("%s Skillboard %s", name, titleSuffix)
	}

	if user.Character != nil && user.Meta != nil && user.QueueSummary != nil {
		meta.Description = fmt.Sprintf(
			"%s is a character in the MMORPG Eve Online. They have amassed %s skillpoints and they're currently training %d skills in their skill queue. Click to learn more about this character",
			user.Character.Name,
			printer.Sprintf("%v", number.Decimal(user.Meta.TotalSP)),
			len(user.QueueSummary.Queue),
		)
	}
	return "meta", meta
}

func (s *Service) userSettingsMeta(ctx context.Context, user *skillz.User) (string, string) {

	name := pluralize.NewClient().Plural(user.Character.Name)
	title := fmt.Sprintf("%s Skillboard Settings %s", name, titleSuffix)

	return "title", title

}

var userSettingsPageTitle = func(name string) (string, string) {
	name = pluralize.NewClient().Plural(name)
	return "title", fmt.Sprintf("%s Settings %s", name, titleSuffix)
}

func (s *Service) userHandler(c buffalo.Context) error {
	var ctx = c.Request().Context()

	userID := c.Param("userID")

	userUUID, err := uuid.FromString(userID)
	if err == nil {
		user, err := s.user.UserByUUID(ctx, userUUID)
		if err != nil {
			if err != nil {
				return c.Error(http.StatusInternalServerError, fmt.Errorf("unable to find user for valid uuid in database"))
			}
		}
		data := render.Data{"userID": user.ID}
		if c.Param("token") != "" {
			data["token"] = c.Param("token")
		}

		return c.Redirect(http.StatusFound, "userPath()", data)
	}

	u, err := s.user.User(ctx, userID)
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return c.Error(http.StatusInternalServerError, err)
	}

	if errors.Is(err, user.ErrUserNotFound) {
		c.Flash().Add("danger", err.Error())
		return c.Redirect(http.StatusFound, "rootPath()")
	}

	settings := u.Settings
	spew.Dump(settings)
	if settings != nil {
		if settings.Visibility == skillz.VisibilityPrivate {
			sessionUser := c.Data()[keyAuthenticatedUser]
			if sessionUser == nil {
				s.flashDanger(c, "User Not Found")
				return c.Redirect(http.StatusFound, "rootPath()")
			}

			authenticatedUser := sessionUser.(*skillz.User)
			if authenticatedUser.ID != userID {
				s.flashDanger(c, "User Not Found")
				return c.Redirect(http.StatusFound, "rootPath()")
			}
		}

		if settings.Visibility == skillz.VisibilityToken {
			token := c.Param("token")
			if token != "" {

				if token != "" && token != settings.VisibilityToken {
					s.flashDanger(c, "User Not Found")
					return c.Redirect(http.StatusFound, "rootPath()")
				}
			} else {
				sessionUser := c.Data()[keyAuthenticatedUser]
				if sessionUser == nil {
					s.flashDanger(c, "User Not Found")
					return c.Redirect(http.StatusFound, "rootPath()")
				}

				authenticatedUser, ok := sessionUser.(*skillz.User)
				if !ok {
					s.flashDanger(c, "User Not Found")
					return c.Redirect(http.StatusFound, "rootPath()")
				}

				if authenticatedUser.ID != userID {
					s.flashDanger(c, "User Not Found")
					return c.Redirect(http.StatusFound, "rootPath()")
				}
			}

		}
	}

	if u.IsNew {
		c.Set("title", fmt.Sprintf("Welcome to Skillboard.Evie %s", titleSuffix))
		return c.Render(http.StatusOK, s.renderer.HTML("user/welcome.plush.html"))
	}

	u, err = s.user.LoadUserAll(ctx, userID)
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return c.Error(http.StatusInternalServerError, err)
	}
	if u.Character != nil {
		c.Set(s.userPageMeta(ctx, u))
	}

	c.Set("user", u)

	return c.Render(http.StatusOK, s.renderer.HTML("user/index.plush.html", "user/layout.plush.html"))
}

func (s *Service) userSettingsHandler(c buffalo.Context) error {

	var ctx = c.Request().Context()

	user := c.Data()[keyAuthenticatedUser].(*skillz.User)
	if user == nil {
		s.logger.Debug("user is missing from session")
		return c.Redirect(http.StatusFound, "rootPath()")
	}

	c.Set(s.userSettingsMeta(ctx, user))
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
		return c.Redirect(http.StatusFound, "rootPath()")
	}

	settings := new(skillz.UserSettings)
	settings.VisibilityToken = user.Settings.VisibilityToken

	err := c.Bind(settings)
	if err != nil {
		s.flashDanger(c, "failed to process form. Please try again")
		return c.Redirect(http.StatusFound, "usersSettingsPath()")
	}

	if !settings.Visibility.Valid() {
		s.flashDanger(c, "Invalid value for Visibilty. Please try again")
		return c.Redirect(http.StatusFound, "usersSettingsPath()")
	}

	err = s.user.CreateUserSettings(ctx, user.ID, settings)
	if err != nil {
		return err
	}

	user.Settings = settings
	c.Set(userSettingsPageTitle(user.Character.Name))

	s.flashSuccess(c, "Settings updated successfully")
	return c.Redirect(http.StatusFound, "usersSettingsPath()")

}

func (s *Service) deleteUserSettingsHandler(c buffalo.Context) error {
	var r = c.Request()
	var ctx = r.Context()
	var form = r.Form

	spew.Dump(form)

	if !form.Has("confirmed") {
		return c.Redirect(http.StatusFound, "usersSettingsPath()")
	}

	if form.Get("confirmed") != "true" {
		return c.Redirect(http.StatusFound, "usersSettingsPath()")
	}

	user := c.Data()[keyAuthenticatedUser].(*skillz.User)
	if user == nil {
		return c.Redirect(http.StatusFound, "rootPath()")
	}

	err := s.user.DeleteUser(ctx, user)
	if err != nil {
		s.flashDanger(c, "Failed to delete user. Please contact the maintainer")
		return c.Redirect(http.StatusFound, "usersSettingsPath()")
	}

	c.Session().Clear()
	s.flashSuccess(c, "You account has successfully been deleted. Goodbye :-(")
	return c.Redirect(http.StatusFound, "rootPath()")

}
