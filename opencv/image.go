package opencv

//#cgo pkg-config: --libs-only-L opencv
//#cgo CFLAGS: -Wno-error=unused-function
//#cgo LDFLAGS: -lopencv_imgproc -lopencv_core -lopencv_highgui
//#include "opencv.h"
import "C"

import (
	"errors"
	"image"
	"image/color"
	"runtime"
	"unsafe"
)

// Wrap IplImage

type Image struct {
	iplImage *C.IplImage
}

func newImage(iplImage *C.IplImage) *Image {
	image := &Image{iplImage}
	runtime.SetFinalizer(image, func(img *Image) { img.Release() })
	return image
}

func Decode(b []byte) (img *Image, err error) {
	defer recoverWithError(&err)

	cvMat := C.cvCreateMatHeader(1, C.int(len(b)), C.CV_8UC1)
	C.cvSetData(unsafe.Pointer(cvMat), unsafe.Pointer(&b[0]), C.int(len(b)))

	iplImage := C.cvDecodeImage(cvMat, C.CV_LOAD_IMAGE_COLOR)
	C.cvReleaseMat(&cvMat)

	return newImage(iplImage), err
}

func (img *Image) Resize(width, height int) (resizedImg *Image, err error) {
	defer recoverWithError(&err)

	resizedImg = img.cloneResizeTarget(width, height)

	C.cvResize(
		unsafe.Pointer(img.iplImage),
		unsafe.Pointer(resizedImg.iplImage),
		C.int(C.CV_INTER_AREA),
	)

	return
}

func (img *Image) Fit(width, height int) (resizedImg *Image, err error) {
	defer recoverWithError(&err)

	if width <= 0 || height <= 0 {
		resizedImg = img
		return
	}

	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	if srcW <= width && srcH <= height {
		resizedImg = img
		return
	}

	srcAspectRatio := float64(srcW) / float64(srcH)
	maxAspectRatio := float64(width) / float64(height)

	var newW, newH int
	if srcAspectRatio > maxAspectRatio {
		newW = width
		newH = int(float64(newW) / srcAspectRatio)
	} else {
		newH = height
		newW = int(float64(newH) * srcAspectRatio)
	}

	resizedImg, err = img.Resize(newW, newH)

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
	if img.iplImage.nChannels == 1 {
		return color.GrayModel
	} else {
		return color.RGBAModel
	}
}

func (img *Image) At(x, y int) color.Color {
	scalar := C.cvGet2D(unsafe.Pointer(img.iplImage), C.int(y), C.int(x))

	if img.iplImage.nChannels == 1 {
		return color.Gray{uint8(scalar.val[0])}
	} else {
		// Convert OpenCV's BGR representation to RGB, which image.Image expects
		return color.RGBA{
			uint8(scalar.val[2]), uint8(scalar.val[1]), uint8(scalar.val[0]), 255,
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
	return newImage(newIpl)
}

func recoverWithError(err *error) {
	if r := recover(); r != nil {
		if _, ok := r.(runtime.Error); ok {
			panic(r)
		}
		switch r.(type) {
		case string:
			*err = errors.New(r.(string))
		default:
			*err = r.(error)
		}
	}
}
