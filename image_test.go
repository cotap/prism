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

func init() {
	jpg, _ := ioutil.ReadFile("./testdata/lenna.jpg")
	lennaJPG = jpg
	png, _ := ioutil.ReadFile("./testdata/lenna.png")
	lennaPNG = png
}

func TestDecodeJPEG(t *testing.T) {
	img, _ := Decode(bytes.NewBuffer(lennaJPG))
	assert.Equal(t, "968a15332343a2794fe7b55f65bd02635e173aad", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
}

func TestDecodePNG(t *testing.T) {
	img, _ := Decode(bytes.NewBuffer(lennaPNG))
	assert.Equal(t, "af1319afa76b41e2ac0f856ccfc267ac262be22c", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
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
