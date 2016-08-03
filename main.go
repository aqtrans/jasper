package main

// Should be able to use font.MeasureString to do what I need:
// https://github.com/golang/freetype/pull/23
// https://github.com/golang/freetype/blob/master/example/drawer/main.go
// https://godoc.org/golang.org/x/image/font

import (
	"image"
	"github.com/golang/freetype"
	"image/draw"
	"os"
	"log"
	"io/ioutil"
	"image/png"
	"math"
	"golang.org/x/image/math/fixed"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"net/http"
	"github.com/dimfeld/httptreemux"
	//"image/color"
	//"fmt"
)

func drawHandler(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(httptreemux.ParamsContextKey).(map[string]string)
	text := params["text"]
	title := "That's a Paddlin'"
	log.Println(text)

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
			Size: fontSize,
			DPI: 72,
			Hinting: font.HintingNone,
		})

	d := &font.Drawer{
		Dst: newimage,
		Src: image.White,
		Face: face,
	}
	d.Dot = fixed.Point26_6{
		X: (fixed.I(originalimage.Bounds().Dx())  - d.MeasureString(title))  / 2,
		Y: fixed.I(originalimage.Bounds().Max.Y - 20 ),
	}
	d.DrawString(title)
	
	//dy := int(math.Ceil(12 * 1.0 * 72 / 72))

	dm := d.MeasureString(text)
	textWidth := dm.Round()
	imageWidth := b.Max.X
	log.Println("textWidth")
	log.Println(textWidth)
	log.Println("imageWidth")
	log.Println(imageWidth)
		
	for textWidth > imageWidth {
		log.Println("Text too long")
		fontSize = fontSize-1.0
		face = truetype.NewFace(myFont, &truetype.Options{
			Size: fontSize,
			DPI: 72,
			Hinting: font.HintingNone,
		})
		d = &font.Drawer{
			Dst: newimage,
			Src: image.White,
			Face: face,
		}		
		dm = d.MeasureString(text)
		textWidth = dm.Round()
		log.Println("textWidth")
		log.Println(textWidth)		
	}

	y := 10 + int(math.Ceil(fontSize*72/72))

	d.Dot = fixed.Point26_6{
		X: (fixed.I(originalimage.Bounds().Dx())  - d.MeasureString(text))  / 2,
		Y: fixed.I(y),
	}
	d.DrawString(text)

	//toimg, _ := os.Create("omg.png")
	//defer toimg.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/png")


	err = png.Encode(w, newimage)
	if err != nil {
		log.Println(err)
	}

}

func main() {
	r := httptreemux.New()
	r.GET("/*text", drawHandler)
	http.Handle("/", r)

	http.ListenAndServe("127.0.0.1:3002", nil)


	//var text = string("OMG YEAH WHAT OMG YEAH WHAT WHAT WHAT WHAT WHAT WHAT WHAT WHAT WHAT WHAT")


	/*
	y += dy
	for _, s := range text {
		d.Dot = fixed.P(10, y)
		d.DrawString(s)
		y += dy
	}
	*/

	// The below works, printing a simple OMG on the image
	// But it does not calculate the width at all
	/* 
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(myFont)
	//c.SetSrc(m)
	c.SetSrc(image.Black)
	c.SetDst(newimage)
	c.SetFontSize(20)
	c.SetClip(originalimage.Bounds())
	imgP := image.Point{0, originalimage.Bounds().Dy()-10}
	pt := freetype.Pt(imgP.X, imgP.Y)
	c.DrawString(text, pt)
	*/

	//c.SetClip(m.Bounds())
	/*
	pt := freetype.Pt(0, b.Dy())
	_, err = c.DrawString("OMGGGG", pt)
	if err != nil {
		log.Fatal(err)
		return
	}
	toimg, _ := os.Create("omg.png")
	defer toimg.Close()
	png.Encode(toimg, newimage)*/

	// http://stackoverflow.com/questions/29105540/aligning-text-in-golang-with-truetype
    // Draw the guidelines.
	/*
    ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
    for rcount := 0; rcount < 4; rcount ++ {
        for i := 0; i < 200; i++ {
            newimage.Set(250*rcount, i, ruler)
        }
    }

    // Truetype stuff
    opts := truetype.Options{}
    opts.Size = 100.0
    face := truetype.NewFace(font, &opts)


    // Calculate the widths and print to image
    for i, x := range(text) {
        awidth, ok := face.GlyphAdvance(rune(x))
        if ok != true {
            log.Println(err)
            return
        }
        iwidthf := int(float64(awidth) / 64)
        //fmt.Printf("%+v\n", iwidthf)

        pt := freetype.Pt(i*250+(125-iwidthf/2), 128)
        c.DrawString(string(x), pt)
        log.Println(string(x))
        //fmt.Printf("%+v\n", awidth)
    }
	*/



}