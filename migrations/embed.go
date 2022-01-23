package migrations

import (
	"embed"
	"io/fs"

	"github.com/gobuffalo/buffalo"
)

//go:embed *
var files embed.FS

func NewFS() fs.FS {
	return buffalo.NewFS(files, ".")
}
