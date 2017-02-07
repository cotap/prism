prism
=====

Fast image processing in Go, powered by OpenCV and libjpeg-turbo.

## Dependencies

- opencv 2.4.x
- libturbo-jpeg 1.4+

## Example

```go
import(
  "os"
  "image/jpeg"

  "github.com/cotap/prism"
)

func main() {
  f, _ := os.Open("example.jpg")

  img, _ := prism.Decode(f)
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
