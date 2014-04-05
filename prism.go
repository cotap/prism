package prism

//#cgo pkg-config: --libs-only-L opencv
//#cgo CFLAGS: -Wno-error=unused-function
//#cgo LDFLAGS: -lopencv_imgproc -lopencv_core -lopencv_highgui
//#include "opencv.h"
import "C"

import (
	"errors"
	"io"
	"io/ioutil"
	"runtime"
	"unsafe"
)

func Decode(r io.Reader) (img *Image, err error) {
	defer recoverWithError(&err)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	cvMat := C.cvCreateMatHeader(1, C.int(len(b)), C.CV_8UC1)
	C.cvSetData(unsafe.Pointer(cvMat), unsafe.Pointer(&b[0]), C.int(len(b)))

	img = newImage(C.cvDecodeImage(cvMat, C.CV_LOAD_IMAGE_COLOR))
	C.cvReleaseMat(&cvMat)

	return
}

func Resize(img *Image, width, height int) (resizedImg *Image, err error) {
	defer recoverWithError(&err)

	resizedImg = img.cloneResizeTarget(width, height)

	C.cvResize(
		unsafe.Pointer(img.iplImage),
		unsafe.Pointer(resizedImg.iplImage),
		C.int(C.CV_INTER_AREA),
	)

	return
}

func Fit(img *Image, width, height int) (resizedImg *Image, err error) {
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

	resizedImg, err = Resize(img, newW, newH)

	return
}

func Rotate90(img *Image) (rotatedImg *Image, err error) {
	defer recoverWithError(&err)

	rotatedImg = img.cloneResizeTarget(img.Bounds().Size().Y, img.Bounds().Size().X)

	C.cvTranspose(
		unsafe.Pointer(img.iplImage),
		unsafe.Pointer(rotatedImg.iplImage),
	)

	return
}

func Rotate180(img *Image) (rotatedImg *Image, err error) {
	rotatedImg, err = Flip(img, -1)
	return
}

func Rotate270(img *Image) (rotatedImg *Image, err error) {
	defer recoverWithError(&err)

	rotatedImg = img.cloneResizeTarget(img.Bounds().Size().Y, img.Bounds().Size().X)

	C.cvTranspose(
		unsafe.Pointer(img.iplImage),
		unsafe.Pointer(rotatedImg.iplImage),
	)

	C.cvFlip(
		unsafe.Pointer(rotatedImg.iplImage),
		unsafe.Pointer(rotatedImg.iplImage),
		C.int(1),
	)

	return
}

func FlipH(img *Image) (flippedImg *Image, err error) {
	flippedImg, err = Flip(img, 1)
	return
}

func FlipV(img *Image) (flippedImg *Image, err error) {
	flippedImg, err = Flip(img, 0)
	return
}

func Flip(img *Image, axis int) (flippedImg *Image, err error) {
	defer recoverWithError(&err)

	flippedImg = img.cloneTarget()

	C.cvFlip(
		unsafe.Pointer(img.iplImage),
		unsafe.Pointer(flippedImg.iplImage),
		C.int(axis),
	)

	return
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
