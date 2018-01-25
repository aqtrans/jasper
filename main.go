package main

// Should be able to use font.MeasureString to do what I need:
// https://github.com/golang/freetype/pull/23
// https://github.com/golang/freetype/blob/master/example/drawer/main.go
// https://godoc.org/golang.org/x/image/font

import (
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"strconv"

	"github.com/dimfeld/httptreemux"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/muesli/cache2go"
	_ "github.com/tevjef/go-runtime-metrics/expvar"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var cache *cache2go.CacheTable

func drawHandler(w http.ResponseWriter, r *http.Request) {
	ptext := httptreemux.ContextParams(r.Context())["text"]
	// Add a question mark to the end of given text
	text := ptext + "?"
	title := "That's a Paddlin'"
	log.Println(text)

	// Try and find image in cache
	cached, cacheErr := cache.Value(ptext)
	if cacheErr == nil {
		log.Println("cached image found: ", cached.Key(), cached.AccessCount())

		w.WriteHeader(http.StatusOK)

		w.Header().Set("Content-Type", "image/png")

		encodeErr := png.Encode(w, cached.Data().(image.Image))
		if encodeErr != nil {
			log.Println(encodeErr)
		}
		return
	}

	reader, err := os.Open("tap.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	originalimage, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	b := originalimage.Bounds()
	newimage := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newimage, newimage.Bounds(), originalimage, image.ZP, draw.Src)

	fontfile, err := ioutil.ReadFile("./DejaVuSansCondensed-Bold.ttf")
	if err != nil {
		log.Fatal(err)
		return
	}
	myFont, err := freetype.ParseFont(fontfile)
	if err != nil {
		log.Fatal(err)
		return
	}

	// First draw That's a Paddlin' at the bottom
	fontSize := 70.0
	face := truetype.NewFace(myFont, &truetype.Options{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	d := &font.Drawer{
		Dst:  newimage,
		Src:  image.White,
		Face: face,
	}
	d.Dot = fixed.Point26_6{
		X: (fixed.I(originalimage.Bounds().Dx()) - d.MeasureString(title)) / 2,
		Y: fixed.I(originalimage.Bounds().Max.Y - 20),
	}
	d.DrawString(title)

	// Now we setup and draw the given text
	dm := d.MeasureString(text)
	textWidth := dm.Round()
	imageWidth := b.Max.X

	// If the width of the text is wider than the image,
	// we loop through shrinking the font size until the text fits
	for textWidth > imageWidth {
		log.Println("Text too long")
		fontSize = fontSize - 1.0
		face = truetype.NewFace(myFont, &truetype.Options{
			Size:    fontSize,
			DPI:     72,
			Hinting: font.HintingNone,
		})
		d = &font.Drawer{
			Dst:  newimage,
			Src:  image.White,
			Face: face,
		}
		dm = d.MeasureString(text)
		textWidth = dm.Round()
		log.Println("textWidth")
		log.Println(textWidth)
	}

	y := 10 + int(math.Ceil(fontSize*72/72))

	d.Dot = fixed.Point26_6{
		X: (fixed.I(originalimage.Bounds().Dx()) - d.MeasureString(text)) / 2,
		Y: fixed.I(y),
	}
	d.DrawString(text)

	// Add element to cache
	cache.Add(ptext, 24*time.Hour, newimage)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/png")

	err = png.Encode(w, newimage)
	if err != nil {
		log.Println(err)
	}

}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.URL.Path)
	if r.URL.Path == "/favicon.ico" {
		serveContent(w, r, "/favicon.ico")
		return
	} else if r.URL.Path == "/favicon.png" {
		serveContent(w, r, "/favicon.png")
		return
	} else {
		http.NotFound(w, r)
		return
	}

}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.URL.Path)
	if r.URL.Path == "/robots.txt" {
		serveContent(w, r, "/robots.txt")
		return
	}
	http.NotFound(w, r)
}

func serveContent(w http.ResponseWriter, r *http.Request, file string) {
	f, err := http.Dir("./").Open(file)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	content := io.ReadSeeker(f)
	http.ServeContent(w, r, file, time.Now(), content)
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	serveContent(w, r, "/tap.png")
	return
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	count := strconv.Itoa(cache.Count())
	w.Write([]byte("<html><body>"))
	w.Write([]byte("<p>Count: " + count + "</p>"))
	w.Write([]byte("<table><thead>"))
	w.Write([]byte("<tr><th>Title</th><th>Access Count</th></tr></thead><tbody>"))
	mostAccessed := cache.MostAccessed(100)
	for _, v := range mostAccessed {
		w.Write([]byte("<tr>"))
		w.Write([]byte("<td>" + v.Key().(string) + "</td>"))
		w.Write([]byte("<td>" + strconv.FormatInt(v.AccessCount(), 10) + "</td>"))
		w.Write([]byte("</tr>"))
	}
	w.Write([]byte("</tbody></table>"))
	w.Write([]byte("</body></html>"))
}

func main() {
	// Initialize the cache
	cache = cache2go.Cache("tap")

	r := httptreemux.NewContextMux()
	r.GET("/_stats", statsHandler)
	r.GET("/*text", drawHandler)
	r.GET("/", indexHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/favicon.png", faviconHandler)
	http.HandleFunc("/robots.txt", http.NotFound)
	http.HandleFunc("/blog", http.NotFound)
	http.HandleFunc("/wp-login.php", http.NotFound)
	http.Handle("/", r)

	log.Println("Now listening on 127.0.0.1:8002")
	http.ListenAndServe("127.0.0.1:8002", nil)
}
