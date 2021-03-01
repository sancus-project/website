package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"

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

func Router() *chi.Mux {
	r := chi.NewRouter()
	r.MethodFunc("GET", "/", HandleIndex)

	return r
}
