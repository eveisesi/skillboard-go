package web

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

func (s *Service) robotsHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, s.renderer.String(`User-agent: *
	Disallow: /
	`))
}
