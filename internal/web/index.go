package web

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/pkg/errors"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

func (s *Service) indexHandler(c buffalo.Context) error {

	r := c.Request()
	ctx := r.Context()

	if s.activeCache {
		p := r.URL.Path
		page, err := s.cache.Page(ctx, p)
		if err == nil && len(page) > 0 {
			w := c.Response()
			_, err = w.Write(page)
			if err != nil {
				s.logger.WithError(err).Error("failed to write bytes to writer")
				return c.Error(http.StatusInternalServerError, err)
			}
			return nil
		}
	}

	c.Set(defaultTitle())

	highlighted, users, err := s.user.Recent(ctx)
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.Wrap(err, "failed to fetch users for homepage"))
	}

	c.Set("highlighted", highlighted)
	c.Set("users", users)

	return s.render(c, http.StatusOK, "text/html", "home/index.plush.html")

}

func (s *Service) render(c buffalo.Context, status int, contentType, template string) error {

	var ctx = c.Request().Context()
	var p = c.Request().URL.Path

	return c.Render(http.StatusOK, s.renderer.Func(contentType, func(w io.Writer, d render.Data) error {

		bb := &bytes.Buffer{}

		err := s.renderer.HTML(template).Render(bb, d)
		if err != nil {
			return err
		}

		m := minify.New()
		m.AddFunc(contentType, html.Minify)

		minified, err := m.Bytes(contentType, bb.Bytes())
		if err != nil {
			s.logger.WithError(err).Error("failed to minifiy document")
			_, err = w.Write(bb.Bytes())
			if err != nil {
				return c.Error(http.StatusInternalServerError, err)
			}
		}

		defer func() {
			err := s.cache.SetPage(ctx, p, minified, time.Minute*30)
			if err != nil {
				s.logger.WithError(err).WithField("p", p).Error("failed to cache page in redis")
			}
		}()

		_, err = w.Write(minified)
		return err

	}))

}
