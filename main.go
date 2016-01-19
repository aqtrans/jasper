package main

import (
	"image"
	"github.com/golang/freetype"
	"image/draw"
	"os"
	"log"
	"io/ioutil"
	"image/png"
)

func main() {
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
	draw.Draw(newimage, newimage.Bounds(), originalimage, b.Min, draw.Src)

	fontfile, err := ioutil.ReadFile("/usr/share/fonts/TTF/LiberationSerif-Regular.ttf")
	if err != nil {
		log.Fatal(err)
		return
	}
	font, err := freetype.ParseFont(fontfile)
	if err != nil {
		log.Fatal(err)
		return
	}	

	c := freetype.NewContext()
	c.SetFont(font)
	//c.SetSrc(m)
	c.SetSrc(image.Black)
	c.SetDst(newimage)
	c.SetFontSize(100)
	c.SetClip(originalimage.Bounds())
	//c.SetClip(m.Bounds())
	pt := freetype.Pt(0, b.Dy())
	_, err = c.DrawString("OMGGGG", pt)
	if err != nil {
		log.Fatal(err)
		return
	}
	toimg, _ := os.Create("omg.png")
	defer toimg.Close()

	png.Encode(toimg, newimage)

}