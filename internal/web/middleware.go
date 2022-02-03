package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func (s *Service) setBaseDomain(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Set("baseDomain", s.baseDomain)
		return next(c)
	}
}

func (s *Service) setCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uuidInf := c.Session().Get(keyAuthenticatedUserID); uuidInf != nil {
			userID, ok := uuidInf.(uuid.UUID)
			if !ok {
				c.Flash().Add("error", "unable to verify authenticated user session")
				c.Session().Delete(keyAuthenticatedUserID)
				return c.Redirect(http.StatusFound, "rootPath()")
			}

			var ctx = c.Request().Context()

			user, err := s.user.User(ctx, userID, user.UserCharacterRel)
			if err != nil {
				c.Flash().Add("danger", fmt.Sprintf("%s is not a valid user id", userID))
				c.Session().Delete(keyAuthenticatedUserID)
				return c.Redirect(http.StatusFound, "rootPath()")
			}

			c.Set(keyAuthenticatedUser, user)
		}
		return next(c)
	}
}

func (s *Service) authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if userID := c.Session().Get(keyAuthenticatedUserID); userID == nil {
			c.Flash().Add("danger", "You must be logged in to see that page.")
			return c.Redirect(http.StatusFound, "rootPath()")
		}
		return next(c)
	}
}

func (s *Service) monitoring(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx := s.newrelic.StartTransaction(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		tx.AddAttribute("environment", s.env.String())
		if r.Referer() != "" {
			tx.AddAttribute("referer", r.Referer())
		}
		defer tx.End()

		tx.SetWebRequestHTTP(r)
		w = tx.SetWebResponse(w)

		r = newrelic.RequestWithTransactionContext(r, tx)

		next.ServeHTTP(w, r)

		p := r.URL.Path
		ignorable := []string{
			"/robots.txt/",
			"/assets/",
		}
		for _, ignore := range ignorable {
			if strings.Contains(p, ignore) {
				tx.Ignore()
			}
		}

	})
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
