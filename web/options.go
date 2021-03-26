package web

import (
	"go.sancus.dev/file2go/html"
)

type Router struct {
	HashifyAssets bool

	html html.Collection
}
