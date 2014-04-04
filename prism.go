
package main

import(
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/cotap/prism/opencv"
)

func main() {
	if (len(os.Args) < 2) {
		fmt.Println("Usage: prism IMAGE")
		os.Exit(1)
	}

	// load original
  name := os.Args[1]
  b, err := ioutil.ReadFile(name)
  if err != nil {
  	panic(err)
  }

	// resize
	img, err := opencv.Decode(b)
	if (err != nil) {
		panic(err)
	}

	width, _ := strconv.Atoi(os.Args[2])
	height, _ := strconv.Atoi(os.Args[3])
	img, err = img.Fit(width, height)
	if (err != nil) {
		panic(err)
	}

	// write resized
	w, _ := os.Create("resized.jpg")
	defer w.Close()
	jpeg.Encode(w, img, &jpeg.Options{ Quality: 90 })

	// open preview
	err = exec.Command("open", "-a", "/Applications/Preview.app", w.Name()).Run()
	if err != nil {
		log.Fatal(err)
	}
}
