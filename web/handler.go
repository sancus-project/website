package web

import (
	"log"
	"net/http"
	"os"
	"time"

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
