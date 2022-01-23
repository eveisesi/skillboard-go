package web

import (
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
)

func (s *Service) SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uuidInf := c.Session().Get(keyAuthenticatedUserID); uuidInf != nil {
			userID, ok := uuidInf.(uuid.UUID)
			if !ok {
				c.Flash().Add("error", "unable to verify authenticated user session")
				c.Session().Delete(keyAuthenticatedUserID)
				return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
			}

			var ctx = c.Request().Context()

			user, err := s.user.User(ctx, userID, user.UserCharacterRel)
			if err != nil {
				c.Flash().Add("danger", fmt.Sprintf("%s is not a valid user id", userID))
				c.Session().Delete(keyAuthenticatedUserID)
				return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
			}

			c.Set(keyAuthenticatedUser, user)
		}
		return next(c)
	}
}

func (s *Service) Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if userID := c.Session().Get(keyAuthenticatedUserID); userID == nil {
			c.Flash().Add("da", "You must be logged in to see that page.")
			return c.Redirect(http.StatusTemporaryRedirect, "rootPath()")
		}
		return next(c)
	}
}

// func (s *Service) PageCacher(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// c := &cacheWriter{ResponseWriter: w}
// 		fmt.Printf("URL: %s\n", r.URL.String())
// 		h.ServeHTTP(w, r)
// 	})
// }

// type cacheWriter struct {
// 	data []byte
// 	http.ResponseWriter
// }

// func (c *cacheWriter) Write(d []byte) (int, error) {
// 	c.data = append(make([]byte, len(d)), d...)
// 	return c.ResponseWriter.Write(d)
// }
