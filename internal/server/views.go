package server

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/eveisesi/skillz/internal/graphql"
	"github.com/go-chi/chi/v5"
	"github.com/tdewolff/minify/v2/minify"
	"github.com/volatiletech/null"
)

var homepageTemplates = template.Must(template.New("home.page.html").Funcs(funcMap).ParseFiles(
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

var funcMap = template.FuncMap{
	"nullUint": func(i null.Uint) uint {
		return i.Uint
	},
	"nullFloat64": func(i null.Float64) float64 {
		return i.Float64
	},
	"fmtDate": func(t time.Time) string {
		return t.Format("2006-01-02")
	},
	"uintToInt64": func(i uint) int64 {
		return int64(i)
	},
	"comma": humanize.Comma,
	"fmtSecStatus": func(f float64) string {
		return humanize.FtoaWithDigits(f, 2)
	},
}

var characterTemplates = template.Must(template.New("character.page.html").Funcs(funcMap).ParseFiles(
	"./views/character.page.html",
	"./views/base.layout.html",
	"./views/footer.partial.html",
))

type characterViewData struct {
	Character     *graphql.Character `json:"character"`
	Skills        *graphql.Skills    `json:"skills"`
	GroupedSkills [][]*graphql.Skills
	Implants      []*graphql.Implant  `json:"implants"`
	Attributes    *graphql.Attributes `json:"attributes"`
	Queue         []*graphql.Queue    `json:"queue"`
}

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

	skillboard, gqlErr := s.graphql.Skillboard(ctx, characterID)
	if gqlErr != nil {
		fmt.Println(gqlErr)
		return
	}

	viewData := &characterViewData{
		Character:  skillboard.User.Character,
		Skills:     skillboard.User.Skills,
		Implants:   skillboard.User.Implants,
		Attributes: skillboard.User.Attributes,
		Queue:      skillboard.User.Queue,
	}

	buf := make([]byte, 0, 4096)
	b := bytes.NewBuffer(buf)

	err = characterTemplates.Execute(b, viewData)
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

// func (s *server) genGroupedSkills(ctx context.Context, skills []*graphql.Skill) ([][]*graphql.Skill, error) {

// 	groupMap := make(map[uint][]*graphql.Skill)
// 	for _, skill := range skills {
// 		if skill.Info.Gr
// 	}

// }
