package web

import (
	"go.sancus.dev/file2go/html"
	"go.sancus.dev/middleware/goget"
)

type Router struct {
	HashifyAssets  bool
	GoImportConfig string

	goget *goget.Config
	html  html.Collection
}
