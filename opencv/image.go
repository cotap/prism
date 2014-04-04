package opencv

//#cgo pkg-config: --libs-only-L opencv
//#cgo CFLAGS: -Wno-error=unused-function
//#cgo LDFLAGS: -lopencv_imgproc -lopencv_core -lopencv_highgui
//#include "opencv.h"
import "C"

import(
	"errors"
	goimage "image"
	"image/color"
	"runtime"
	"unsafe"
)

// Wrap IplImage

type Image struct {
	config   goimage.Config
	iplImage *C.IplImage
}

func Decode(b []byte) (*Image, error) {
	cvMat := C.cvCreateMatHeader(1, C.int(len(b)), C.CV_8UC1)
	C.cvSetData(unsafe.Pointer(cvMat), unsafe.Pointer(&b[0]), C.int(len(b)))

	iplImage := C.cvDecodeImage(cvMat, C.CV_LOAD_IMAGE_COLOR)
	C.cvReleaseMat(&cvMat)

	return imageFromIplImage(iplImage)
}

func (img *Image) Resize(width, height int) (*Image, error) {
	var resizedIpl *C.IplImage

	size := C.CvSize { width: C.int(width), height: C.int(height) }
	resizedIpl = C.cvCreateImage(size , img.iplImage.depth, img.iplImage.nChannels)
	C.cvResize(unsafe.Pointer(img.iplImage), unsafe.Pointer(resizedIpl), C.int(C.CV_INTER_CUBIC))

	resizedImg, err := imageFromIplImage(resizedIpl)
	if err != nil {
		return nil, err
	}

	return resizedImg, nil
}

func (img *Image) Fit(width, height int) (*Image, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("dimensions less than 0")
	}

	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	if srcW <= width && srcH <= height {
		return img, nil
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

func (img *Image) Release() {
	if img.iplImage != nil {
		C.cvReleaseImage(&img.iplImage)
		img.iplImage = nil
	}
}


// image.Image interface

func (img *Image) ColorModel() color.Model {
	return img.config.ColorModel
}

func (img *Image) At(x, y int) color.Color {
	scalar := C.cvGet2D(unsafe.Pointer(img.iplImage), C.int(y), C.int(x))

	if img.iplImage.nChannels == 1 {
		return color.Gray { uint8(scalar.val[0]) }
	} else {
		// Convert OpenCV's BGR representation to RGB, which image.Image expects
		return color.RGBA {
			uint8(scalar.val[2]), uint8(scalar.val[1]), uint8(scalar.val[0]), 255,
		}
	}
}

func (img *Image) Bounds() goimage.Rectangle {
	return goimage.Rect(0, 0, img.config.Width, img.config.Height)
}


// init from *C.IplImage

func imageFromIplImage(iplImage *C.IplImage) (*Image, error) {
	var model color.Model

	switch(iplImage.nChannels) {
		case 1: model = color.GrayModel
		case 3: model = color.RGBAModel
		case 4: model = color.RGBAModel
		default:
			return nil, errors.New("Unsupported image type - number of channels")
	}

	size := C.cvGetSize(unsafe.Pointer(iplImage))
	config := goimage.Config {
		ColorModel: model,
		Width: int(size.width),
		Height: int(size.height),
	}

	image := &Image { config, iplImage }
	runtime.SetFinalizer(image, func(img *Image) { img.Release() })

	return image, nil
}
