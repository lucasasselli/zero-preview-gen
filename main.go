package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
)

func main() {

	const fps = 60
	const seconds = 5
	const layersDir = "layers"
	const tempPath = "/tmp/zero"
	const depth = 1

	// get arguments
	args := os.Args

	if len(args) != 2 {
		log.Fatal("ERROR: argument format not correct!")
	}

	workPath := args[1]

	// search and delete temp
	cmd := exec.Command("rm", "-fr", tempPath)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	// Create temp folder
	cmd = exec.Command("mkdir", tempPath) 
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	files, _ := ioutil.ReadDir(workPath + "/" + layersDir)
	for _, f := range files {
		fmt.Println("Opening file ", f.Name())
	}

	var id int

	for i := 0; i < fps*seconds; i++ {
		fmt.Printf("%d/%d\n", i, seconds*fps)

		K := 2 * math.Pi / (fps * seconds)

		img := image.NewRGBA(image.Rect(0, 0, 2000, 2000))
		for _, f := range files {

			name := f.Name()

			var pos int
			var z int

			fmt.Sscanf(name, "%d_%d_%d.png", &pos, &id, &z)

			offX := int(float64(z) * depth * math.Cos(K*float64(i)))
			offY := int(float64(z) * depth * math.Sin(K*float64(i)))

			file, _ := os.Open(workPath + "/" + layersDir + "/" + f.Name())
			defer file.Close()
			layer, _, _ := image.Decode(file)

			draw.Draw(img, img.Bounds(), layer, image.Point{offX, offY}, draw.Over)

		}
		croppedImg, _ := cutter.Crop(img, cutter.Config{
			Width:  1800,
			Height: 1800,
			Mode:   cutter.Centered,
		})
		croppedImg = resize.Resize(500, 500, croppedImg, resize.MitchellNetravali)

		outName := fmt.Sprintf("out_%03d.png", i)

		outFile, _ := os.Create(tempPath + "/" + outName)
		defer outFile.Close()
		png.Encode(outFile, croppedImg)
	}

	videoName := fmt.Sprintf("%05d.mp4", id)
	cmd = exec.Command("ffmpeg", "-r", strconv.Itoa(fps), "-f", "image2", "-s", "500X500", "-i", tempPath+"/out_%03d.png", "-vcodec", "libx264", "-crf", "15", "-pix_fmt", "yuv420p", workPath+"/"+videoName)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
