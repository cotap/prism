package prism

//#cgo pkg-config: --libs-only-L opencv
//#cgo CFLAGS: -Wno-error=unused-function
//#cgo LDFLAGS: -lopencv_imgproc -lopencv_core -lopencv_highgui
//#include "opencv.h"
import "C"

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"runtime"
	"unsafe"

	"github.com/rwcarlsen/goexif/exif"
)

// Wrap IplImage

type Image struct {
	iplImage *C.IplImage
	exif     *exif.Exif
}

func newImage(iplImage *C.IplImage, meta *exif.Exif) *Image {
	image := &Image{iplImage, meta}
	runtime.SetFinalizer(image, func(img *Image) { img.Release() })
	return image
}

func Decode(r io.Reader) (img *Image, err error) {
	defer recoverWithError(&err)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	cvMat := C.cvCreateMatHeader(1, C.int(len(b)), C.CV_8UC1)
	C.cvSetData(unsafe.Pointer(cvMat), unsafe.Pointer(&b[0]), C.int(len(b)))

	meta, _ := exif.Decode(bytes.NewReader(b))

	iplImage := C.cvDecodeImage(cvMat, C.CV_LOAD_IMAGE_UNCHANGED)
	if iplImage == nil {
		err = errors.New("Bad image")
		return
	}

	img = newImage(iplImage, meta)
	C.cvReleaseMat(&cvMat)

	return
}

func (img *Image) Release() {
	if img.iplImage != nil {
		C.cvReleaseImage(&img.iplImage)
		img.iplImage = nil
	}
}

// image.Image interface

func (img *Image) ColorModel() color.Model {
	switch img.iplImage.nChannels * img.iplImage.depth {
	case 8:
		return color.GrayModel
	case 16:
		return color.Gray16Model
	case 48, 64:
		return color.NRGBA64Model
	case 24, 32:
		fallthrough
	default:
		return color.NRGBAModel
	}
}

func (img *Image) At(x, y int) color.Color {
	scalar := C.cvGet2D(unsafe.Pointer(img.iplImage), C.int(y), C.int(x))

	// Convert OpenCV's BGRA representation to RGBA, which image.Image expects
	switch img.ColorModel() {
	case color.GrayModel:
		return color.Gray{uint8(scalar.val[0])}
	case color.Gray16Model:
		return color.Gray16{uint16(scalar.val[0])}
	case color.NRGBA64Model:
		alpha := ^uint16(0)
		if img.iplImage.nChannels == 4 {
			alpha = uint16(scalar.val[3])
		}

		return color.NRGBA64{
			uint16(scalar.val[2]),
			uint16(scalar.val[1]),
			uint16(scalar.val[0]),
			alpha,
		}
	case color.NRGBAModel:
		fallthrough
	default:
		alpha := ^uint8(0)
		if img.iplImage.nChannels == 4 {
			alpha = uint8(scalar.val[3])
		}

		return color.NRGBA{
			uint8(scalar.val[2]),
			uint8(scalar.val[1]),
			uint8(scalar.val[0]),
			alpha,
		}
	}
}

func (img *Image) Bounds() image.Rectangle {
	size := C.cvGetSize(unsafe.Pointer(img.iplImage))
	return image.Rect(0, 0, int(size.width), int(size.height))
}

// create target image

func (img *Image) cloneTarget() *Image {
	return img.cloneResizeTarget(img.Bounds().Size().X, img.Bounds().Size().Y)
}

// create target image with new size, but same color depth and channels

func (img *Image) cloneResizeTarget(width, height int) *Image {
	size := C.CvSize{width: C.int(width), height: C.int(height)}
	newIpl := C.cvCreateImage(size, img.iplImage.depth, img.iplImage.nChannels)
	return newImage(newIpl, img.exif)
}
