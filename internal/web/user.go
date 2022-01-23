package web

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/gertd/go-pluralize"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gofrs/uuid"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
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
	r := c.Request()
	ctx := r.Context()
	p := r.URL.Path

	page, err := s.cache.Page(ctx, p)
	if err == nil && len(page) > 0 {
		w := c.Response()
		w.Write(page)
		return nil
	}

	userID, err := uuid.FromString(c.Param("userID"))
	if err != nil {
		c.Flash().Add("error", "invalid id supplied for userID")
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

	if u.IsNew {
		return c.Render(http.StatusOK, s.renderer.HTML("user/new.plush.html"))
	}

	u, err = s.user.LoadUserAll(ctx, userID)
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return c.Error(http.StatusInternalServerError, err)
	}
	c.Set(defaultTitle())
	if u.Character != nil {
		c.Set(userPageTitle(u.Character.Name))
	}

	c.Set("user", u)

	return c.Render(http.StatusOK, s.renderer.Func("text/html", func(w io.Writer, d render.Data) error {
		bb := &bytes.Buffer{}

		err := s.renderer.HTML("user/index.plush.html").Render(bb, d)
		if err != nil {
			return err
		}

		m := minify.New()
		m.AddFunc("text/html", html.Minify)

		minified, err := m.Bytes("text/html", bb.Bytes())
		if err != nil {
			s.logger.WithError(err).Error("failed to minifiy document")
			_, err = w.Write(bb.Bytes())
			if err != nil {
				return c.Error(http.StatusInternalServerError, err)
			}
		}

		defer s.cache.SetPage(ctx, p, minified, time.Minute*5)

		_, err = w.Write(minified)
		return err

	}))
}

func (s *Service) userSettingsHandler(c buffalo.Context) error {

	authedUser := c.Data()[keyAuthenticatedUser].(*skillz.User)
	if authedUser == nil {
		return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
	}

	c.Set(userSettingsPageTitle(authedUser.Character.Name))

	return c.Render(http.StatusOK, s.renderer.HTML("user/settings.plush.html"))
}

func (s *Service) postUserSettingsHandler(c buffalo.Context) error {

	settings := new(skillz.UserSettings)

	err := c.Bind(settings)
	if err != nil {
		return err
	}

	spew.Dump(settings)
	return nil

}
