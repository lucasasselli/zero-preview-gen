package main

import (
	"fmt"
	"github.com/andybons/gogif"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/gif"
	_ "image/png"
	"io/ioutil"
	"math"
	"os"
)

func main() {

	const fps = 100
	const seconds = 3
	const path = "./layers/"
	const depth = 0.5

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		fmt.Println("Opening file ", f.Name())
	}

	outGif := &gif.GIF{}

	for i := 0; i < fps*seconds; i++ {
		fmt.Printf("%d/%d\n", i, seconds*fps)

		K := 2 * math.Pi / (fps * seconds)

		img := image.NewRGBA(image.Rect(0, 0, 500, 500))
		for _, f := range files {

			name := f.Name()

			var pos int
			var id int
			var z int

			fmt.Sscanf(name, "%d_%d_%d.png", &pos, &id, &z)

			fmt.Println(z)

			offX := int(float64(z) * depth * math.Cos(K*float64(i)))
			offY := int(float64(z) * depth * math.Sin(K*float64(i)))

			file, _ := os.Open(path + f.Name())
			defer file.Close()
			layer, _, _ := image.Decode(file)
			layer = resize.Resize(500, 500, layer, resize.Lanczos3)
			draw.Draw(img, img.Bounds(), layer, image.Point{offX, offY}, draw.Over)
		}

		palettedImg := image.NewPaletted(img.Bounds(), nil)
		quantizer := gogif.MedianCutQuantizer{NumColor: 256}
		quantizer.Quantize(palettedImg, img.Bounds(), img, image.ZP)

		outGif.Image = append(outGif.Image, palettedImg)
		outGif.Delay = append(outGif.Delay, 1)
	}

	// save to out.gif
	f, _ := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, outGif)
}
