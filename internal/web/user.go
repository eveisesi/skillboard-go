package web

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
)

func (s *Service) userHandler(c buffalo.Context) error {
	id := c.Param("userID")
	return c.Render(http.StatusOK, s.renderer.String(fmt.Sprintf("Hello World with ID of %s", id)))
}
