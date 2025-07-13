package web

import (
	"embed"
	"io/fs"

	"github.com/akatranlp/sentinel/utils"
)

//go:embed all:dist
var assetEmbedFS embed.FS

var AssetFS fs.FS

//go:embed dist_templates/*.tmpl.html
var templatesEmbedFS embed.FS

var TemplateFS fs.FS

func init() {
	AssetFS = assetEmbedFS

	TemplateFS = utils.Must(fs.Sub(templatesEmbedFS, "dist_templates"))
}
