package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/cotap/prism"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: prism IMAGE WIDTH HEIGHT")
		os.Exit(1)
	}
	width, _ := strconv.Atoi(os.Args[2])
	height, _ := strconv.Atoi(os.Args[3])

	// load original
	name := os.Args[1]
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, format, _ := image.DecodeConfig(io.Reader(f))
	f.Seek(0, 0)

	img, err := prism.Decode(io.Reader(f))
	if err != nil {
		panic(err)
	}

	reoriented, err := prism.Reorient(img)
	if err != nil {
		panic(err)
	}

	resized, err := prism.Fit(reoriented, width, height)
	if err != nil {
		panic(err)
	}

	// write resized
	w, _ := os.Create("resized." + format)
	defer w.Close()

	switch format {
	case "jpeg":
		err = jpeg.Encode(w, resized, &jpeg.Options{Quality: 85})
	case "png":
		err = png.Encode(w, resized)
	}

	if err != nil {
		panic(err)
	}

	// open preview
	err = exec.Command("open", "-a", "/Applications/Preview.app", w.Name()).Run()
	if err != nil {
		log.Fatal(err)
	}
}
