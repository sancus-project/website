package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.sancus.dev/middleware/goget"

	"github.com/sancus-project/website/assets"
	"github.com/sancus-project/website/html"
)

func (c *Router) Compile() error {

	// bind assets to html templates
	h, err := html.Files.Clone()
	if err != nil {
		return err
	}
	h.Funcs(assets.Files.FuncMap(c.HashifyAssets, "File"))
	// compile templates
	if err := h.Parse(); err != nil {
		return err
	}

	c.html = h

	// go-import middleware
	c.goget = &goget.Config{
		Filename: c.GoImportConfig,
	}
	if err := c.goget.Load(); err != nil {
		return err
	}
	return nil
}

func (c Router) Reload() error {
	if err := c.goget.Reload(); err != nil {
		return err
	}
	return nil
}

func (c *Router) Handler() http.Handler {
	// and compose the router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(assets.Middleware(c.HashifyAssets))
	r.Use(c.goget.NewMiddleware(
		goget.OnlyGoGet{}, // only renderer go-import for ?go-get=1 requests
		goget.RedirectToDoc{}, // if the package exists but no ?go-get=1 was given
	))
	r.Use(html.Middleware(c.html))
	r.MethodFunc("GET", "/", HandleIndex)

	return r
}
