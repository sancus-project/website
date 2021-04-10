package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
	return nil
}

func (c Router) Reload() error {
	return nil
}

func (c *Router) Handler() http.Handler {
	// and compose the router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(assets.Middleware(c.HashifyAssets))
	r.Use(html.Middleware(c.html))
	r.MethodFunc("GET", "/", HandleIndex)

	return r
}
