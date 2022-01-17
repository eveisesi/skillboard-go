package templates

import (
	"embed"
	"io/fs"

	"github.com/gobuffalo/buffalo"
)

//go:embeded * */*
var files embed.FS

func FS() fs.FS {
	return buffalo.NewFS(files, "internal/templates")
}
