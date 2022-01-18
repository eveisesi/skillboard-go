package web

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
)

func (s *Service) userHandler(c buffalo.Context) error {
	id := c.Param("userID")
	return c.Render(http.StatusOK, s.renderer.String(fmt.Sprintf("Hello World with ID of %s", id)))
}

func (s *Service) userStatusHandler(c buffalo.Context) error {
	id, err := uuid.FromString(c.Param("userID"))
	if err != nil {
		return err
	}

	user, err := s.user.User(c.Request().Context(), id, nil)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch user by id")
		return c.Render(http.StatusInternalServerError, s.renderer.JSON(struct {
			Message string `json:"message"`
		}{
			Message: "failed to fetch user by id",
		}))
	}

	return c.Render(http.StatusOK, s.renderer.JSON(struct {
		IsNew bool `json:"is_new"`
	}{
		IsNew: user.IsNew,
	}))

}
