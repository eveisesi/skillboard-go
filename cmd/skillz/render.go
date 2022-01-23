package main

import (
	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/skill"
	"github.com/eveisesi/skillz/internal/templates"
	"github.com/eveisesi/skillz/public"
	"github.com/gobuffalo/buffalo/render"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

func renderer() *render.Engine {

	printer := message.NewPrinter(language.English)

	return render.New(render.Options{
		HTMLLayout:  "layout.plush.html",
		TemplatesFS: templates.FS(),
		AssetsFS:    public.FS(),
		Helpers: render.Helpers{
			"formatNum": func(v interface{}) string {
				return printer.Sprintf("%v", number.Decimal(v))
			},
			"addUint": func(v ...uint) uint {
				i := uint(0)
				for _, vv := range v {
					i += vv
				}

				return i
			},
			"nextPos": func(queue []*skillz.CharacterSkillQueue) *skillz.CharacterSkillQueue {
				var out *skillz.CharacterSkillQueue
				if len(queue) == 0 {
					return out
				}
				for _, position := range queue {
					if out == nil {
						out = position
						continue
					}
					if position.QueuePosition < out.QueuePosition {
						out = position
					}
				}
				return out
			},
			"int":               func(v uint) int { return int(v) },
			"hasSkillCount":     skill.PlushHelperHasSkillCount,
			"missingSkillCount": skill.PlushHelperMissingSkill,
			"levelVSkillCount":  skill.PlushHelperLevelVSkillCount,
			"levelVSPTotal":     skill.PlushHelperLevelVSPTotal,
			"possibleSP":        skill.PlushHelperPossibleSPByRank,
		},
	})
}
