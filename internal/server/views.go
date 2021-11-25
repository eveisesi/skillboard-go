package server

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tdewolff/minify/v2/minify"
)

var homepageTemplates = template.Must(template.ParseFiles(
	"./views/home.page.html",
	"./views/base.layout.html",
	"./views/footer.partial.html",
))

func (s *server) handleRenderHomepage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleRenderHomepage")

	err := homepageTemplates.Execute(w, nil)
	if err != nil {
		s.logger.WithError(err).Error("failed to execute template")
		return
	}

}

var characterTemplates = template.Must(template.ParseFiles(
	"./views/character.page.html",
	"./views/base.layout.html",
	"./views/footer.partial.html",
))

func (s *server) handleRenderCharacterPage(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	var url = new(url.URL)
	*url = *r.URL
	url.RawQuery = ""
	url.RawFragment = ""

	// fmt.Println(url.String())
	// page, err := s.cache.Page(ctx, url.String())
	// if err != nil {
	// 	panic(err)
	// }

	// if page != "" {
	// 	_, err = w.Write([]byte(page))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	return
	// }

	characterID, err := strconv.ParseUint(chi.URLParam(r, "userUUID"), 10, 64)
	if err != nil {
		return
	}

	data := s.graphql.Skillboard(ctx, characterID)

	buf := make([]byte, 0, 4096)
	b := bytes.NewBuffer(buf)

	err = characterTemplates.Execute(b, data)
	if err != nil {
		s.logger.WithError(err).Error("failed to execute template")
		return
	}

	minified, err := minify.HTML(b.String())
	if err != nil {
		panic(err)
	}

	// err = s.cache.SetPage(ctx, url.String(), minified, time.Minute*10)
	// if err != nil {
	// 	panic(err)
	// }

	_, err = w.Write([]byte(minified))
	if err != nil {
		panic(err)
	}

}
