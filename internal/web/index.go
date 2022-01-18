package web

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

var defaultTitle = func() (string, string) {
	return "title", "Skillboard || Eve Is ESI || A Third Party Eve Online App"
}

func (s *Service) indexHandler(c buffalo.Context) error {
	c.Set(defaultTitle())

	return c.Render(http.StatusOK, s.renderer.Func("text/html", func(w io.Writer, d render.Data) error {
		bb := &bytes.Buffer{}

		err := s.renderer.HTML("home/index.plush.html").Render(bb, d)
		if err != nil {
			return err
		}

		m := minify.New()
		m.AddFunc("text/html", html.Minify)

		minified, _ := m.Bytes("text/html", bb.Bytes())

		_, err = w.Write(minified)
		return err

	}))
}
