// +build dev

package assets

import (
	"net/http"
)

// Assets contains project assets.
var Assets http.FileSystem = http.Dir("assets")

func FaviconICO(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/favicon.ico")
	return
}

func FaviconPNG(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/favicon.png")
	return
}
