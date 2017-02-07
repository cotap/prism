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

func (img *Image) Resize(width, height int) (err error) {
	img.m.Lock()
	defer img.m.Unlock()
	defer recoverWithError(&err)

	resizedIplImg := allocTarget(img.iplImage, width, height)

	C.cvResize(
		unsafe.Pointer(img.iplImage),
		unsafe.Pointer(resizedIplImg),
		C.int(C.CV_INTER_AREA),
	)

	C.cvReleaseImage(&img.iplImage)
	img.iplImage = resizedIplImg

	return nil
}

func (img *Image) Reorient() (err error) {
	if img.exif == nil {
		return nil
	}

	orientation, err := img.exif.Get(exif.Orientation)
	if err != nil {
		return err
	}

	orientationValue, err := orientation.Int(0)
	if err != nil {
		return err
	}

	switch orientationValue {
	case 2:
		err = img.FlipH()
	case 3:
		err = img.Rotate180()
	case 4:
		err = img.FlipV()
	case 5:
		if err = img.Rotate90(); err == nil {
			err = img.FlipH()
		}
	case 6:
		err = img.Rotate90()
	case 7:
		if err = img.Rotate270(); err == nil {
			err = img.FlipH()
		}
	case 8:
		err = img.Rotate270()
	}

	return err
}

func (img *Image) Fit(width, height int) (err error) {
	if width <= 0 && height <= 0 {
		return
	}

	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	if srcW <= width && srcH <= height {
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

	return img.Resize(newW, newH)
}

func (img *Image) Rotate90() (err error) {
	img.m.Lock()
	defer img.m.Unlock()
	defer recoverWithError(&err)

	newIplImg := allocTarget(img.iplImage, img.Bounds().Size().Y, img.Bounds().Size().X)

	C.cvTranspose(
		unsafe.Pointer(img.iplImage),
		unsafe.Pointer(newIplImg),
	)

	C.cvReleaseImage(&img.iplImage)
	img.iplImage = newIplImg

	return flip(newIplImg, 1)
}

func (img *Image) Rotate180() error {
	img.m.Lock()
	defer img.m.Unlock()
	return flip(img.iplImage, -1)
}

func (img *Image) Rotate270() (err error) {
	img.m.Lock()
	defer img.m.Unlock()
	defer recoverWithError(&err)

	newIplImg := allocTarget(img.iplImage, img.Bounds().Size().Y, img.Bounds().Size().X)

	C.cvTranspose(
		unsafe.Pointer(img.iplImage),
		unsafe.Pointer(newIplImg),
	)

	C.cvReleaseImage(&img.iplImage)
	img.iplImage = newIplImg

	return flip(newIplImg, 0)
}

func (img *Image) FlipH() error {
	img.m.Lock()
	defer img.m.Unlock()
	return flip(img.iplImage, 1)
}

func (img *Image) FlipV() error {
	img.m.Lock()
	defer img.m.Unlock()
	return flip(img.iplImage, 0)
}

func flip(iplImage *C.IplImage, axis int) (err error) {
	defer recoverWithError(&err)

	C.cvFlip(
		unsafe.Pointer(iplImage),
		unsafe.Pointer(iplImage),
		C.int(axis),
	)

	return nil
}

// create target image with new size, but same color depth and channels
func allocTarget(iplImage *C.IplImage, width, height int) *C.IplImage {
	size := C.CvSize{width: C.int(width), height: C.int(height)}
	return C.cvCreateImage(size, iplImage.depth, iplImage.nChannels)
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
