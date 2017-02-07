package prism

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var lenna *Image
var gradRgb48 *Image
var gradRgba64 *Image
var gray *Image

func init() {
	lenna = testImg("lenna.png")
}

func TestEncodeJPEG85(t *testing.T) {
	var enc bytes.Buffer
	EncodeJPEG(bufio.NewWriter(&enc), lenna, 85)

	assert.Equal(t, "0aa716b8d3d263486d12ca097e88700b7d5034f1", fmt.Sprintf("%x", sha1.Sum(enc.Bytes())))
}

func TestEncodeJPEG100(t *testing.T) {
	var enc bytes.Buffer
	EncodeJPEG(bufio.NewWriter(&enc), lenna, 100)

	assert.Equal(t, "bb4a165387fd90cd245f2df06c9501f01144689d", fmt.Sprintf("%x", sha1.Sum(enc.Bytes())))
}

func TestEncodeJPEGRGB48(t *testing.T) {
	var enc bytes.Buffer
	EncodeJPEG(bufio.NewWriter(&enc), testImg("rgb48.png"), 85)

	assert.Equal(t, "fd6e4ab757bef9d733ae33ed69d19c3be0e61d49", fmt.Sprintf("%x", sha1.Sum(enc.Bytes())))
}

func TestEncodeJPEGRGBA64(t *testing.T) {
	var enc bytes.Buffer
	EncodeJPEG(bufio.NewWriter(&enc), testImg("rgba64.png"), 85)

	assert.Equal(t, "f1818c85e6615725cef12622a6d7a3ac731ae5b9", fmt.Sprintf("%x", sha1.Sum(enc.Bytes())))
}

func TestEncodeJPEGGray(t *testing.T) {
	var enc bytes.Buffer
	EncodeJPEG(bufio.NewWriter(&enc), testImg("gray.jpg"), 85)

	assert.Equal(t, "71e9a2ed07b5251abbbf14687618a0178ae1bc2c", fmt.Sprintf("%x", sha1.Sum(enc.Bytes())))
}

func TestEncodePNG0(t *testing.T) {
	var enc bytes.Buffer
	EncodePNG(bufio.NewWriter(&enc), lenna, 0)

	assert.Equal(t, "ee5097475f6e5178da8fd340952131b40cabbfc6", fmt.Sprintf("%x", sha1.Sum(enc.Bytes())))
}

func TestEncodePNG9(t *testing.T) {
	var enc bytes.Buffer
	EncodePNG(bufio.NewWriter(&enc), lenna, 9)

	assert.Equal(t, "338afb0ca3fb268c0fc4220aee3723968d6d68cd", fmt.Sprintf("%x", sha1.Sum(enc.Bytes())))
}

func BenchmarkEncodeJPEG85(b *testing.B) {
	for n := 0; n < b.N; n++ {
		EncodeJPEG(ioutil.Discard, lenna, 85)
	}
}

func BenchmarkEncodeJPEG100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		EncodeJPEG(ioutil.Discard, lenna, 100)
	}
}

func BenchmarkEncodePNG0(b *testing.B) {
	for n := 0; n < b.N; n++ {
		EncodePNG(ioutil.Discard, lenna, 0)
	}
}

func BenchmarkEncodePNG4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		EncodePNG(ioutil.Discard, lenna, 4)
	}
}

func BenchmarkEncodePNG9(b *testing.B) {
	for n := 0; n < b.N; n++ {
		EncodePNG(ioutil.Discard, lenna, 9)
	}
}
