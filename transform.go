package prism

//#cgo pkg-config: --libs-only-L opencv libturbojpeg
//#cgo CFLAGS: -O3 -Wno-error=unused-function
//#cgo LDFLAGS: -lopencv_imgproc -lopencv_core -lopencv_highgui -lturbojpeg
//#include "prism.h"
import "C"

import (
	"errors"
	"runtime"
	"unsafe"

	"github.com/rwcarlsen/goexif/exif"
)

func Resize(img *Image, width, height int) (resizedImg *Image, err error) {
	defer recoverWithError(&err)

	resizedImg = img.cloneResizeTarget(width, height)

	C.cvResize(
		unsafe.Pointer(img.iplImage),
		unsafe.Pointer(resizedImg.iplImage),
		C.int(C.CV_INTER_AREA),
	)

	return resizedImg, nil
}

func Reorient(img *Image) (*Image, error) {
	if img.exif != nil {
		reorientedImg, err := reorientByExif(img)
		if err != nil {
			return nil, err
		}
		return reorientedImg, nil
	}
	return img, nil
}

func Fit(img *Image, width, height int) (resizedImg *Image, err error) {
	defer recoverWithError(&err)

	if width <= 0 && height <= 0 {
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

	if width == 0 {
		width = srcW
	}

	if height == 0 {
		height = srcH
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

	return Resize(img, newW, newH)
}

func reorientByExif(img *Image) (*Image, error) {
	var err error

	orientation, err := img.exif.Get(exif.Orientation)
	if err != nil {
		return nil, err
	}

	orientationValue, err := orientation.Int(0)
	if err != nil {
		return nil, err
	}

	switch orientationValue {
	case 2:
		img, err = FlipH(img)
	case 3:
		img, err = Rotate180(img)
	case 4:
		img, err = FlipV(img)
	case 5:
		img, err = Rotate90(img)
		if err == nil {
			img, err = FlipH(img)
		}
	case 6:
		img, err = Rotate90(img)
	case 7:
		img, err = Rotate270(img)
		if err == nil {
			img, err = FlipH(img)
		}
	case 8:
		img, err = Rotate270(img)
	}

	return img, err
}

func Rotate90(img *Image) (_ *Image, err error) {
	defer recoverWithError(&err)

	rotatedImg := img.cloneResizeTarget(img.Bounds().Size().Y, img.Bounds().Size().X)

	C.cvTranspose(
		unsafe.Pointer(img.iplImage),
		unsafe.Pointer(rotatedImg.iplImage),
	)

	return rotatedImg, flip(rotatedImg, rotatedImg, 1)
}

func Rotate180(img *Image) (*Image, error) {
	flippedImg := img.cloneTarget()
	return flippedImg, flip(img, flippedImg, -1)
}

func Rotate270(img *Image) (_ *Image, err error) {
	defer recoverWithError(&err)

	rotatedImg := img.cloneResizeTarget(img.Bounds().Size().Y, img.Bounds().Size().X)

	C.cvTranspose(
		unsafe.Pointer(img.iplImage),
		unsafe.Pointer(rotatedImg.iplImage),
	)

	return rotatedImg, flip(rotatedImg, rotatedImg, 0)
}

func FlipH(img *Image) (*Image, error) {
	flippedImg := img.cloneTarget()
	return flippedImg, flip(img, flippedImg, 1)
}

func FlipV(img *Image) (*Image, error) {
	flippedImg := img.cloneTarget()
	return flippedImg, flip(img, flippedImg, 0)
}

func flip(img *Image, flippedImg *Image, axis int) (err error) {
	defer recoverWithError(&err)

	C.cvFlip(
		unsafe.Pointer(img.iplImage),
		unsafe.Pointer(flippedImg.iplImage),
		C.int(axis),
	)

	return nil
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
