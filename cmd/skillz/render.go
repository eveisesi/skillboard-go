package main

import (
	"github.com/eveisesi/skillz/internal/templates"
	"github.com/eveisesi/skillz/public"
	"github.com/gobuffalo/buffalo/render"
)

func renderer() *render.Engine {
	return render.New(render.Options{
		HTMLLayout:  "layout.plush.html",
		TemplatesFS: templates.FS(),
		AssetsFS:    public.FS(),
	})
}
