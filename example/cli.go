package main

import (
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/cotap/prism"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: prism IMAGE")
		os.Exit(1)
	}

	// load original
	name := os.Args[1]
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	width, _ := strconv.Atoi(os.Args[2])
	height, _ := strconv.Atoi(os.Args[3])

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
	w, _ := os.Create("resized.jpg")
	defer w.Close()
	jpeg.Encode(w, resized, &jpeg.Options{Quality: 90})

	// open preview
	err = exec.Command("open", "-a", "/Applications/Preview.app", w.Name()).Run()
	if err != nil {
		log.Fatal(err)
	}
}
