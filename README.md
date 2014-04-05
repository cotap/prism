prism
=====

Fast image processing in Go, powered by OpenCV.

## Example

```go
import(
  "os"
  "image/jpeg"

  "github.com/cotap/prism"
)

func main() {
  f, _ := os.Open("example.jpg")

  var img prism.Image
  img, _ = prism.Decode(f)
  img, _ = prism.FlipV(img)
  img, _ = prism.Fit(img, 500, 500)

  w, _ := os.Create("resized.jpg")
  defer w.Close()

  // prism.Image implements image.Image
  jpeg.Encode(w, img, &jpeg.Options{Quality: 90})
}
```
