package html

import (
	"context"
	"log"
	"net/http"

	"github.com/amery/file2go/html"
	"github.com/go-chi/render"
)

type ctxKey int

const collectionCtxKey ctxKey = 0

func Middleware(h html.Collection) func(http.Handler) http.Handler {
	key := collectionCtxKey

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), key, h))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func View(r *http.Request, name string, data interface{}) render.Renderer {
	key := collectionCtxKey

	if h, ok := r.Context().Value(key).(html.Collection); ok {
		if v, err := h.View(name, data); err == nil {
			return v
		} else {
			log.Fatal(err)
		}
	}

	return nil
}
