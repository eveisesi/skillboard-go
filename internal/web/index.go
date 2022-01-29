package web

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

func (s *Service) indexHandler(c buffalo.Context) error {

	r := c.Request()
	ctx := r.Context()

	highlighted, users, err := s.user.Recent(ctx)
	if err != nil {
		return c.Error(http.StatusInternalServerError, errors.Wrap(err, "failed to fetch users for homepage"))
	}

	c.Set("highlighted", highlighted)
	c.Set("users", users)

	return c.Render(http.StatusOK, s.renderer.HTML("home/index.plush.html", "home/layout.plush.html"))
}
