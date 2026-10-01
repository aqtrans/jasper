package main

// Should be able to use font.MeasureString to do what I need:
// https://github.com/golang/freetype/pull/23
// https://github.com/golang/freetype/blob/master/example/drawer/main.go
// https://godoc.org/golang.org/x/image/font

import (
	"embed"
	"errors"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"net/http"
	"regexp"
	"time"

	//"git.jba.io/go/jasper/vfs"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"

	//"github.com/shurcooL/httpfs/vfsutil"
	"github.com/go-chi/httprate"
	_ "github.com/tevjef/go-runtime-metrics/expvar"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Go 1.16 embed:
//
//go:embed assets
var assetsfs embed.FS

func sanitizeInput(s string) (string, error) {
	re := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
	if re.MatchString(s) {
		return "", errors.New("unable to sanitize text:" + s)
	}
	return re.ReplaceAllString(s, ""), nil
}

func drawHandler(w http.ResponseWriter, r *http.Request) {
	ptext := httptreemux.ContextParams(r.Context())["text"]

	// Do not allow long words
	if len(ptext) > 200 {
		http.Error(w, "Unable to handle this request", http.StatusBadRequest)
		return
	}

	// Sanitize input content
	sanitizedText, err := sanitizeInput(ptext)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Add a question mark to the end of given text
	text := sanitizedText + "?"
	title := "That's a Paddlin'"
	//log.Println(text)

	reader, err := assetsfs.Open("assets/tap.png")
	//reader, err := vfs.VFS.Open("tap.png")
	if err != nil {
		http.Error(w, "Unable to handle this request", http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer reader.Close()

	originalimage, _, err := image.Decode(reader)
	if err != nil {
		http.Error(w, "Unable to handle this request", http.StatusBadRequest)
		log.Println(err)
		return
	}
	b := originalimage.Bounds()
	newimage := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newimage, newimage.Bounds(), originalimage, image.Point{0, 0}, draw.Src)

	fontFile, err := assetsfs.ReadFile("assets/impact.ttf")
	//fontFile, err := vfsutil.ReadFile(vfs.VFS, "impact.ttf")
	if err != nil {
		http.Error(w, "Unable to handle this request", http.StatusBadRequest)
		log.Println("Error loading impact.ttf", err)
		return
	}
	myFont, err := freetype.ParseFont(fontFile)
	if err != nil {
		http.Error(w, "Unable to handle this request", http.StatusBadRequest)
		log.Println(err)
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
		//log.Println("Text too long")
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
		//log.Println("textWidth")
		//log.Println(textWidth)
	}

	y := 10 + int(math.Ceil(fontSize*72/72))

	d.Dot = fixed.Point26_6{
		X: (fixed.I(originalimage.Bounds().Dx()) - d.MeasureString(text)) / 2,
		Y: fixed.I(y),
	}
	d.DrawString(text)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/png")

	err = png.Encode(w, newimage)
	if err != nil {
		http.Error(w, "Unable to handle this request", http.StatusBadRequest)
		log.Println(err)
		return
	}

}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.URL.Path)
	if r.URL.Path == "/favicon.ico" {
		serveContent(w, "assets/favicon.ico")
		return
	} else if r.URL.Path == "/favicon.png" {
		serveContent(w, "assets/favicon.png")
		return
	} else {
		http.NotFound(w, r)
		return
	}

}

func serveContent(w http.ResponseWriter, file string) {
	assetBytes, err := assetsfs.ReadFile(file)
	if err != nil {
		log.Println("error reading file from assetfs:", file, err)
	}
	_, err = w.Write(assetBytes)
	if err != nil {
		log.Println("error writing HTTP response:", file, err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>That's a Paddlin'</title>
    <meta name="description" content="Generate a paddlin' meme with any text. What deserves a paddlin'?">
    <meta property="og:title" content="That's a Paddlin'">
    <meta property="og:description" content="Generate a paddlin' meme with any text.">
    <meta property="og:type" content="website">
    <meta property="og:url" content="https://thatsapaddl.in">
    <meta property="og:image" content="https://thatsapaddl.in/tap.png">
    <link rel="icon" href="/favicon.ico" type="image/x-icon">
    <link rel="icon" href="/favicon.png" type="image/png">
    <style>
        :root {
            --bg: #1a1f2e;
            --bg-card: #242b3d;
            --text: #e8eaf0;
            --text-muted: #8b92a8;
            --accent: #f5c518;
            --accent-hover: #e0b010;
            --accent-active: #cc9d0e;
            --error: #ff6b6b;
            --radius: 12px;
        }
        * { box-sizing: border-box; }
        body {
            margin: 0;
            padding: 2rem 1rem;
            background: linear-gradient(160deg, #1a1f2e 0%, #1e2a3a 50%, #1a2535 100%);
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            color: var(--text);
        }
        .container {
            width: 100%;
            max-width: 640px;
            text-align: center;
        }
        .headline {
            font-size: clamp(2rem, 5vw, 3rem);
            margin: 0 0 0.25rem;
            font-weight: 800;
            letter-spacing: -0.02em;
        }
        .subtext {
            font-size: 1rem;
            color: var(--text-muted);
            margin: 0 0 2rem;
        }
        .card {
            background: var(--bg-card);
            border-radius: var(--radius);
            padding: 2rem;
            box-shadow: 0 8px 32px rgba(0,0,0,0.3);
        }
        form {
            display: flex;
            flex-direction: column;
            gap: 1rem;
            margin-bottom: 1.5rem;
        }
        label {
            font-size: 0.875rem;
            font-weight: 600;
            color: var(--text-muted);
            text-transform: uppercase;
            letter-spacing: 0.08em;
            text-align: left;
        }
        input[type="text"] {
            width: 100%;
            padding: 0.75rem 1rem;
            font-size: 1rem;
            font-family: inherit;
            background: #1a1f2e;
            color: var(--text);
            border: 2px solid #3a4258;
            border-radius: 8px;
            outline: none;
            transition: border-color 0.2s ease, box-shadow 0.2s ease;
        }
        input[type="text"]:focus {
            border-color: var(--accent);
            box-shadow: 0 0 0 3px rgba(245,197,24,0.15);
        }
        input[type="text"]::placeholder {
            color: var(--text-muted);
        }
        button[type="submit"] {
            padding: 0.75rem 2rem;
            font-size: 1rem;
            font-weight: 700;
            font-family: inherit;
            background: var(--accent);
            color: #1a1f2e;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            transition: background 0.2s ease, transform 0.1s ease;
        }
        button[type="submit"]:hover {
            background: var(--accent-hover);
        }
        button[type="submit"]:active {
            background: var(--accent-active);
            transform: scale(0.97);
        }
        button[type="submit"]:focus-visible {
            outline: 2px solid var(--accent);
            outline-offset: 2px;
        }
        button[type="submit"]:disabled {
            opacity: 0.6;
            cursor: not-allowed;
        }
        .error-msg {
            color: var(--error);
            font-size: 0.875rem;
            margin-top: 0.5rem;
            display: none;
        }
        .error-msg.visible {
            display: block;
        }
        .image-area {
            position: relative;
            border-radius: 8px;
            overflow: hidden;
            background: #111;
        }
        .image-area img {
            width: 100%;
            height: auto;
            display: block;
            transition: opacity 0.3s ease;
        }
        .spinner {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            width: 48px;
            height: 48px;
            border: 4px solid rgba(255,255,255,0.1);
            border-top-color: var(--accent);
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
            display: none;
        }
        .spinner.visible {
            display: block;
        }
        @keyframes spin {
            to { transform: translate(-50%, -50%) rotate(360deg); }
        }
        .copy-btn {
            margin-top: 1rem;
            padding: 0.5rem 1rem;
            font-size: 0.875rem;
            font-weight: 600;
            font-family: inherit;
            background: transparent;
            color: var(--text-muted);
            border: 1px solid #3a4258;
            border-radius: 6px;
            cursor: pointer;
            transition: color 0.2s ease, border-color 0.2s ease;
            display: none;
        }
        .copy-btn.visible {
            display: inline-block;
        }
        .copy-btn:hover {
            color: var(--text);
            border-color: var(--text-muted);
        }
        .copy-btn:focus-visible {
            outline: 2px solid var(--accent);
            outline-offset: 2px;
        }
        footer {
            margin-top: 2rem;
            font-size: 0.8rem;
            color: var(--text-muted);
            opacity: 0.6;
        }
        @media (max-width: 480px) {
            .card { padding: 1.25rem; }
            .headline { font-size: 1.75rem; }
        }
    </style>
</head>
<body>
    <main class="container">
        <h1 class="headline">What deserves a paddlin'?</h1>
        <p class="subtext">Type it. Paddle it. Share it.</p>
        <div class="card">
            <form id="paddle-form" action="/paddle" method="post">
                <label for="paddle">Your text</label>
                <input type="text" id="paddle" name="paddle" placeholder="What deserves a paddlin'?" maxlength="200" autocomplete="off">
                <button type="submit" id="paddle-btn">Paddle</button>
                <p class="error-msg" id="error-msg" role="alert"></p>
            </form>
            <div class="image-area" id="status" aria-busy="false">
                <img src="/tap.png" alt="That's a Paddlin' base image" id="preview-img">
                <div class="spinner" id="spinner" aria-hidden="true"></div>
            </div>
            <button class="copy-btn" id="copy-btn" type="button">Copy Link</button>
        </div>
        <footer>That's a Paddlin' &mdash; thatsapaddl.in</footer>
    </main>
    <script>
    (function() {
        var form = document.getElementById('paddle-form');
        var input = document.getElementById('paddle');
        var btn = document.getElementById('paddle-btn');
        var img = document.getElementById('preview-img');
        var spinner = document.getElementById('spinner');
        var errorMsg = document.getElementById('error-msg');
        var copyBtn = document.getElementById('copy-btn');
        var status = document.getElementById('status');

        function showError(msg) {
            errorMsg.textContent = msg;
            errorMsg.classList.add('visible');
        }

        function clearError() {
            errorMsg.classList.remove('visible');
        }

        function setLoading(loading) {
            spinner.classList.toggle('visible', loading);
            btn.disabled = loading;
            btn.textContent = loading ? 'Paddlin\u2019...' : 'Paddle';
            img.style.opacity = loading ? '0.3' : '1';
            status.setAttribute('aria-busy', loading ? 'true' : 'false');
        }

        function showCopy(url) {
            copyBtn.classList.add('visible');
            copyBtn.dataset.url = url;
        }

        form.addEventListener('submit', async function(e) {
            e.preventDefault();
            clearError();

            var text = input.value.trim();

            if (text.length === 0) {
                showError('Please enter some text to paddle.');
                return;
            }
            if (text.length > 200) {
                showError('Text must be 200 characters or fewer.');
                return;
            }
            if (!/^[a-zA-Z0-9 ]+$/.test(text)) {
                showError('Only letters, numbers, and spaces are allowed.');
                return;
            }

            setLoading(true);
            copyBtn.classList.remove('visible');

            try {
                var formData = new FormData();
                formData.append('paddle', text);
                var res = await fetch('/paddle', {
                    method: 'POST',
                    body: formData
                });

                if (!res.ok) {
                    showError('Something went wrong. Please try again.');
                } else {
                    var blob = await res.blob();
                    var url = URL.createObjectURL(blob);
                    img.src = url;
                    img.alt = text + "? That's a paddlin'!";
                    var canonicalUrl = window.location.origin + '/' + encodeURIComponent(text);
                    showCopy(canonicalUrl);
                }
            } catch (err) {
                showError('Network error. Please try again.');
            } finally {
                setLoading(false);
            }
        });

        copyBtn.addEventListener('click', function() {
            var url = copyBtn.dataset.url;
            if (!url) return;
            navigator.clipboard.writeText(url).then(function() {
                copyBtn.textContent = 'Copied!';
                setTimeout(function() {
                    copyBtn.textContent = 'Copy Link';
                }, 2000);
            }).catch(function() {
                window.prompt('Copy this link:', url);
            });
        });
    })();
    </script>
</body>
</html>`))

	if err != nil {
		log.Println("error writing HTTP response: ", err)
	}
}

func formPost(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.FormValue("paddle"))
	http.Redirect(w, r, "/"+r.FormValue("paddle"), http.StatusSeeOther)
}

func tapDirect(w http.ResponseWriter, r *http.Request) {
	serveContent(w, "assets/tap.png")
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write([]byte("<html><body>"))
	if err != nil {
		log.Println("error writing HTTP response: ", err)
	}
	_, err = w.Write([]byte("<table><thead>"))
	if err != nil {
		log.Println("error writing HTTP response: ", err)
	}
	_, err = w.Write([]byte("<tr><th>Title</th><th>Access Count</th></tr></thead><tbody>"))
	if err != nil {
		log.Println("error writing HTTP response: ", err)
	}
	_, err = w.Write([]byte("</tbody></table>"))
	if err != nil {
		log.Println("error writing HTTP response: ", err)
	}
	_, err = w.Write([]byte("</body></html>"))

	if err != nil {
		log.Println("error writing HTTP response: ", err)
	}
}

func main() {
	limit := httprate.LimitByIP(100, time.Minute)

	r := httptreemux.NewContextMux()
	r.UseHandler(limit)
	r.GET("/_stats", statsHandler)
	r.GET("/*text", drawHandler)
	r.POST("/paddle", formPost)
	r.GET("/tap.png", tapDirect)
	r.GET("/", indexHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/favicon.png", faviconHandler)
	http.HandleFunc("/robots.txt", http.NotFound)
	http.HandleFunc("/blog", http.NotFound)
	http.HandleFunc("/wp-login.php", http.NotFound)
	http.Handle("/", r)

	log.Println("Now listening on 127.0.0.1:8002")
	err := http.ListenAndServe("127.0.0.1:8002", nil)
	if err != nil {
		log.Fatalln("error listening ")
	}
}
