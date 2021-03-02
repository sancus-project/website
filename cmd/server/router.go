package main

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

	err := html.Files.ExecuteTemplate(w, "index", vd)
	if err != nil {
		log.Println(err)
	}
}

func Router(hashify bool) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(assets.Middleware(hashify))
	r.MethodFunc("GET", "/", HandleIndex)

	return r
}
