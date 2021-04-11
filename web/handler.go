package web

import (
	"net/http"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/sancus-project", 302)
}
