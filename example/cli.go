package main

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/cotap/prism"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: prism IMAGE WIDTH HEIGHT")
		os.Exit(1)
	}
	width, _ := strconv.Atoi(os.Args[2])
	height, _ := strconv.Atoi(os.Args[3])

	// load original
	name := os.Args[1]
	b, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}

	_, format, err := image.DecodeConfig(bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}

	img, err := prism.Decode(bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}

	_ = img.Reorient()

	err = img.Fit(width, height)
	if err != nil {
		panic(err)
	}

	// write resized
	w, _ := os.Create("resized." + format)
	defer w.Close()

	switch format {
	case "png":
		err = prism.EncodePNG(w, img, 4)
	case "jpg", "jpeg":
		err = prism.EncodeJPEG(w, img, 85)
	}

	if err != nil {
		panic(err)
	}

	// // open preview
	// err = exec.Command("open", "-a", "/Applications/Preview.app", w.Name()).Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
