package prism

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var lennaJPG []byte
var lennaPNG []byte
var bombJPG []byte
var bombPNG []byte

func init() {
	jpg, _ := ioutil.ReadFile("./testdata/lenna.jpg")
	lennaJPG = jpg
	png, _ := ioutil.ReadFile("./testdata/lenna.png")
	lennaPNG = png
	jpg, _ = ioutil.ReadFile("./testdata/bomb.jpg")
	bombJPG = jpg
	png, _ = ioutil.ReadFile("./testdata/bomb.png")
	bombPNG = png
}

func testImg(name string) *Image {
	b, err := ioutil.ReadFile("./testdata/" + name)
	img, err := Decode(bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}
	return img
}

func TestDecodeJPEG(t *testing.T) {
	img, _ := Decode(bytes.NewBuffer(lennaJPG))
	assert.Equal(t, "968a15332343a2794fe7b55f65bd02635e173aad", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
}

func TestDecodePNG(t *testing.T) {
	img, _ := Decode(bytes.NewBuffer(lennaPNG))
	assert.Equal(t, "af1319afa76b41e2ac0f856ccfc267ac262be22c", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
}

func TestDecodeJPEGBomb(t *testing.T) {
	img, err := Decode(bytes.NewBuffer(bombJPG))
	assert.Nil(t, img)
	assert.NotNil(t, err)
	assert.Equal(t, "Image is too large (possible decompression bomb): 25500 x 25500", err.Error())
}

func TestDecodePNGBomb(t *testing.T) {
	img, err := Decode(bytes.NewBuffer(bombPNG))
	assert.Nil(t, img)
	assert.NotNil(t, err)
	assert.Equal(t, "Image is too large (possible decompression bomb): 25500 x 25500", err.Error())
}

func BenchmarkDecodeJPEG(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buffer := bytes.NewBuffer(lennaJPG)
		img, _ := Decode(buffer)
		img.Release()
	}
}

func BenchmarkDecodePNG(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buffer := bytes.NewBuffer(lennaPNG)
		img, _ := Decode(buffer)
		img.Release()
	}
}
