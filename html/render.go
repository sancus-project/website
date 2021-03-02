package html

import (
	"log"

	"github.com/go-chi/render"
)

func View(name string, data interface{}) render.Renderer {
	view, err := Files.View(name, data)

	if err != nil {
		log.Fatalln(err)
	}

	return view
}
