package main

import (
	"fmt"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/skill"
	"github.com/eveisesi/skillz/public"
	"github.com/eveisesi/skillz/templates"
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
			"int":                    func(v uint) int { return int(v) },
			"hasSkillCount":          skill.PlushHelperHasSkillCount,
			"missingSkillCount":      skill.PlushHelperMissingSkill,
			"levelVSkillCount":       skill.PlushHelperLevelVSkillCount,
			"levelVSPTotal":          skill.PlushHelperLevelVSPTotal,
			"possibleSP":             skill.PlushHelperPossibleSPByRank,
			"activeNav":              activeNav,
			"activeTabPane":          activeTabPane,
			"hasShipCount":           hasShipCount,
			"percentageShipsTrained": percentageShipsTrained,
			"notFlyable":             notFlyable,
		},
	})
}

var tabs = []string{"skills", "queue", "flyable", "implants"}

const activeNavClass = "active"
const activeTabPaneClass = "show active"

func activeNav(t string, settings *skillz.UserSettings) string {
	activeTab := ""

	for _, tab := range tabs {
		if tab == "skills" && !settings.HideSkills {
			activeTab = tab
			break
		} else if tab == "queue" && !settings.HideQueue {
			activeTab = tab
			break
		} else if tab == "flyable" && !settings.HideQueue {
			activeTab = tab
			break
		} else if tab == "implants" && !settings.HideQueue {
			activeTab = tab
			break
		}
	}

	if t == activeTab {
		return activeNavClass
	}

	return ""
}

func activeTabPane(t string, settings *skillz.UserSettings) string {
	activeTab := ""

	for _, tab := range tabs {
		if tab == "skills" && !settings.HideSkills {
			activeTab = tab
			break
		} else if tab == "queue" && !settings.HideQueue {
			activeTab = tab
			break
		} else if tab == "flyable" && !settings.HideQueue {
			activeTab = tab
			break
		} else if tab == "implants" && !settings.HideQueue {
			activeTab = tab
			break
		}
	}

	if t == activeTab {
		return activeTabPaneClass
	}

	return ""
}

func hasShipCount(ships []*skillz.ShipType) int {
	var count = 0
	for _, ship := range ships {
		if ship == nil {
			fmt.Println("ship is nil")
			continue
		}
		if ship.Flyable {
			count++
		}
	}
	return count
}

func percentageShipsTrained(ships []*skillz.ShipType) string {
	var count = 0
	for _, ship := range ships {
		if ship == nil {
			fmt.Println("ship is nil")
			continue
		}
		if ship.Flyable {
			count++
		}
	}

	percentage := float64(count) / float64(len(ships)) * float64(100)
	return fmt.Sprintf("%.2f", percentage)
}

func notFlyable(b bool) string {
	if b {
		return ""
	}
	return "notFlyable"
}
