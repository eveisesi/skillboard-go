package server

import (
	"fmt"
	"html/template"
	"net/http"
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

	fmt.Println("handleRenderCharacterPage")

	err := characterTemplates.Execute(w, nil)
	if err != nil {
		s.logger.WithError(err).Error("failed to execute template")
		return
	}

}
