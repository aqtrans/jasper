// +build !dev

package assets

import (
	"log"
	"net/http"
	"time"

	"github.com/shurcooL/httpfs/vfsutil"
)

func serve(name string, w http.ResponseWriter, r *http.Request) {
	file, err := Assets.Open(name)
	if err != nil {
		log.Println("Error opening", name)
		w.Write([]byte(""))
		return
	}
	http.ServeContent(w, r, name, time.Now(), file)
	return
}

func FaviconICO(w http.ResponseWriter, r *http.Request) {
	serve("favicon.ico", w, r)
	return
}

func FaviconPNG(w http.ResponseWriter, r *http.Request) {
	serve("favicon.png", w, r)
	return
}

func Font() []byte {
	fontFile, err := vfsutil.ReadFile(Assets, "DejaVuSansCondensed-Bold.ttf")
	if err != nil {
		log.Println("Error loading DejaVuSansCondensed-Bold.ttf", err)
		return []byte("")
	}
	return fontFile
}
