package web

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/amery/go-webpack-starter/assets"
	"github.com/amery/go-webpack-starter/html"
)

func Router(hashify bool) *chi.Mux {

	// bind assets to html templates
	h, err := html.Files.Clone()
	if err != nil {
		log.Fatal(err)
	}
	h.Funcs(assets.Files.FuncMap(hashify, "File"))
	// compile templates
	if err := h.Parse(); err != nil {
		log.Fatal(err)
	}

	// and compose the router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(assets.Middleware(hashify))
	r.Use(html.Middleware(h))
	r.MethodFunc("GET", "/", HandleIndex)

	return r
}
