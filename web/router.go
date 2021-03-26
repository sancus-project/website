package web

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/amery/go-webpack-starter/assets"
	"github.com/amery/go-webpack-starter/html"
)

type View struct {
	Pid int
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if d := r.URL.Query().Get("delay"); d != "" {
		if delay, err := time.ParseDuration(d); err == nil {
			time.Sleep(delay)
		}
	}

	w.WriteHeader(http.StatusOK)

	vd := View{
		Pid: os.Getpid(),
	}

	err := html.View(r, "index", vd).Render(w, r)
	if err != nil {
		log.Fatalln(err)
	}
}

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
