# Prism

[![Build Status](https://travis-ci.org/cotap/prism.svg?branch=master)](https://travis-ci.org/cotap/prism)

Fast Go image encoding, decoding, and transformations using OpenCV, libpng, and libjpeg-turbo.

## Dependencies (dynamically linked)

- opencv 2.4.x
- libturbo-jpeg 1.4+

## Example

```go
import(
  "os"

  "github.com/cotap/prism"
)

func main() {
  f, _ := os.Open("example.jpg")

  img, _ := prism.Decode(f)
  defer img.Release() // ensures timely release of C heap memory

  _ = img.FlipV()
  _ = img.Fit(500, 500)

  w, _ := os.Create("resized.jpg")
  defer w.Close()

  // fast encoding with libjpeg-turbo
  prism.EncodeJPEG(w, img, 90)

  // or, since prism.Image implements image.Image:
  //   jpeg.Encode(w, img, &jpeg.Options{Quality: 90})
}
```

## Memory Management

Prism allocates memory for image decoding and processing using the OpenCV
allocator on the C heap. In order to avoid long-term leaks, prism uses Go's
`runtime.SetFinalizer` to release C memory when the Go object is GC'd.

*HOWEVER*, because Go's collector can't see into the C heap, it may not run as
frequently as needed to keep memory usage low. In order to ensure most
efficient memory usage, always call ` func (img *prism.Image) Release()` after
you're done with an image.
