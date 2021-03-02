package assets

import (
	"net/http"
)

func Middleware(hashify bool) func(http.Handler) http.Handler {
	return Files.Middleware(hashify)
}
