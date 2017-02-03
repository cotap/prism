package prism

//#cgo pkg-config: --libs-only-L opencv libturbojpeg
//#cgo CFLAGS: -O3 -Wno-error=unused-function
//#cgo LDFLAGS: -lopencv_imgproc -lopencv_core -lopencv_highgui -lturbojpeg
//#include "prism.h"
import "C"
import (
	"errors"
	"image/color"
	"image/jpeg"
	"io"
	"unsafe"
)

// EncodeJPEG writes the Image img to w in JPEG 4:2:0 baseline format with the
// given quality, from 1 - 100.
func EncodeJPEG(w io.Writer, img *Image, quality int) (err error) {
	defer recoverWithError(&err)

	if img.ColorModel() == color.Gray16Model || img.ColorModel() == color.NRGBA64Model {
		// workaround for lack of bitdepth conversion in OpenCV C API
		err = jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
		return
	}

	result := C.prismEncodeJPEG(img.iplImage, C.int(quality))
	if result == nil {
		err = errors.New("Unable to encode JPEG image")
		return
	}

	// write bytes directly without copying to Go-land
	_, err = w.Write((*[1 << 30]byte)(unsafe.Pointer(result.buffer))[:result.size:result.size])
	C.prismRelease(result)
	return
}

// EncodePNG writes the Image img to w in PNG format with the given
// Zlib compression level, from 0 (none) - 9
func EncodePNG(w io.Writer, img *Image, compression int) (err error) {
	defer recoverWithError(&err)

	result := C.prismEncodePNG(img.iplImage, C.int(compression))
	if result == nil {
		err = errors.New("Unable to encode PNG image")
		return
	}

	// write bytes directly without copying to Go-land
	_, err = w.Write((*[1 << 30]byte)(unsafe.Pointer(result.buffer))[:result.size:result.size])
	C.prismRelease(result)
	return
}
