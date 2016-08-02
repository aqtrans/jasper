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
	//"image/color"
	//"fmt"
)

func main() {
	var text = string("OMG YEAH WHAT OMG YEAH WHAT")
	reader, err := os.Open("avatar.png")
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

	fontfile, err := ioutil.ReadFile("/usr/share/fonts/TTF/LiberationSerif-Regular.ttf")
	if err != nil {
		log.Fatal(err)
		return
	}
	myFont, err := freetype.ParseFont(fontfile)
	if err != nil {
		log.Fatal(err)
		return
	}

	d := &font.Drawer{
		Dst: newimage,
		Src: image.Black,
		Face: truetype.NewFace(myFont, &truetype.Options{
			Size: 20,
			DPI: 72,
			Hinting: font.HintingNone,
		}),
	}

	y := 10 + int(math.Ceil(12*72/72))
	//dy := int(math.Ceil(12 * 1.0 * 72 / 72))
	d.Dot = fixed.Point26_6{
		X: (fixed.I(originalimage.Bounds().Dx())  - d.MeasureString(text))  / 2,
		Y: fixed.I(y),
	}
	d.DrawString(text)

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

	toimg, _ := os.Create("omg.png")
	defer toimg.Close()

	png.Encode(toimg, newimage)

}