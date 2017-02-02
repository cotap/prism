package prism

import (
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var nerd *Image

func init() {
	nerd = testImg("nerd.jpg")
}

func TestResize(t *testing.T) {
	resized, _ := Resize(nerd, 100, 100)

	assert.Equal(t, "e92728b1537082ee4a5e1bfd0ff242b31d6ee782", fmt.Sprintf("%x", sha1.Sum(resized.Bytes())))
	assert.Equal(t, 100, resized.Bounds().Dx())
	assert.Equal(t, 100, resized.Bounds().Dy())
}

func TestReorient1(t *testing.T) {
	reoriented, _ := Reorient(testImg("orientations/orientation-1.jpg"))

	assert.Equal(t, "06d6a0c8dbdc3229b8ea8d1fe257c890baa41088", fmt.Sprintf("%x", sha1.Sum(reoriented.Bytes())))
	assert.Equal(t, 480, reoriented.Bounds().Dx())
	assert.Equal(t, 640, reoriented.Bounds().Dy())
}

func TestReorient2(t *testing.T) {
	reoriented, _ := Reorient(testImg("orientations/orientation-2.jpg"))

	assert.Equal(t, "0b12c9d9c59544298929fb229cb6cdcefd2e8432", fmt.Sprintf("%x", sha1.Sum(reoriented.Bytes())))
	assert.Equal(t, 480, reoriented.Bounds().Dx())
	assert.Equal(t, 640, reoriented.Bounds().Dy())
}

func TestReorient3(t *testing.T) {
	reoriented, _ := Reorient(testImg("orientations/orientation-3.jpg"))

	assert.Equal(t, "3802560db62d5703a8d96c821efd3e839995e681", fmt.Sprintf("%x", sha1.Sum(reoriented.Bytes())))
	assert.Equal(t, 480, reoriented.Bounds().Dx())
	assert.Equal(t, 640, reoriented.Bounds().Dy())
}

func TestReorient4(t *testing.T) {
	reoriented, _ := Reorient(testImg("orientations/orientation-4.jpg"))

	assert.Equal(t, "58de3d130761d7837bc61fad72673d4ef328f18f", fmt.Sprintf("%x", sha1.Sum(reoriented.Bytes())))
	assert.Equal(t, 480, reoriented.Bounds().Dx())
	assert.Equal(t, 640, reoriented.Bounds().Dy())
}

func TestReorient5(t *testing.T) {
	reoriented, _ := Reorient(testImg("orientations/orientation-5.jpg"))

	assert.Equal(t, "101f3dc5b9834dbf00db0f0f77f71233bb6defc5", fmt.Sprintf("%x", sha1.Sum(reoriented.Bytes())))
	assert.Equal(t, 480, reoriented.Bounds().Dx())
	assert.Equal(t, 640, reoriented.Bounds().Dy())
}

func TestReorient6(t *testing.T) {
	reoriented, _ := Reorient(testImg("orientations/orientation-6.jpg"))

	assert.Equal(t, "0abcaf927976b94134156e6c0385ea3262a66457", fmt.Sprintf("%x", sha1.Sum(reoriented.Bytes())))
	assert.Equal(t, 480, reoriented.Bounds().Dx())
	assert.Equal(t, 640, reoriented.Bounds().Dy())
}

func TestReorient7(t *testing.T) {
	reoriented, _ := Reorient(testImg("orientations/orientation-7.jpg"))

	assert.Equal(t, "bcf54553b6b7a0157caf25f56147256c75298706", fmt.Sprintf("%x", sha1.Sum(reoriented.Bytes())))
	assert.Equal(t, 480, reoriented.Bounds().Dx())
	assert.Equal(t, 640, reoriented.Bounds().Dy())
}

func TestReorient8(t *testing.T) {
	reoriented, _ := Reorient(testImg("orientations/orientation-8.jpg"))

	assert.Equal(t, "22e656fb927de87e9081c14d34375bb960a9444b", fmt.Sprintf("%x", sha1.Sum(reoriented.Bytes())))
	assert.Equal(t, 480, reoriented.Bounds().Dx())
	assert.Equal(t, 640, reoriented.Bounds().Dy())
}

func TestRotate90(t *testing.T) {
	rotated, _ := Rotate90(nerd)

	assert.Equal(t, "f2f2b083811cbe8090201683d860c30dbec8757a", fmt.Sprintf("%x", sha1.Sum(rotated.Bytes())))
	assert.Equal(t, 653, rotated.Bounds().Dx())
	assert.Equal(t, 400, rotated.Bounds().Dy())
}

func TestRotate180(t *testing.T) {
	rotated, _ := Rotate180(nerd)

	assert.Equal(t, "be49e7e35b61db3ceb3e52603ded7a1e17cd8dff", fmt.Sprintf("%x", sha1.Sum(rotated.Bytes())))
	assert.Equal(t, 400, rotated.Bounds().Dx())
	assert.Equal(t, 653, rotated.Bounds().Dy())
}

func TestRotate270(t *testing.T) {
	rotated, _ := Rotate270(nerd)

	assert.Equal(t, "68dc4bd0104936f533bdd3fab7883e9b0711677f", fmt.Sprintf("%x", sha1.Sum(rotated.Bytes())))
	assert.Equal(t, 653, rotated.Bounds().Dx())
	assert.Equal(t, 400, rotated.Bounds().Dy())
}

func TestFlipH(t *testing.T) {
	flipped, _ := FlipH(nerd)

	assert.Equal(t, "f21ebf2fcda3cfd0b1512eacd9f8858c2b6d505b", fmt.Sprintf("%x", sha1.Sum(flipped.Bytes())))
	assert.Equal(t, 400, flipped.Bounds().Dx())
	assert.Equal(t, 653, flipped.Bounds().Dy())
}

func TestFlipV(t *testing.T) {
	flipped, _ := FlipV(nerd)

	assert.Equal(t, "f19d9ce9121772e4882e3115579332af50eaca32", fmt.Sprintf("%x", sha1.Sum(flipped.Bytes())))
	assert.Equal(t, 400, flipped.Bounds().Dx())
	assert.Equal(t, 653, flipped.Bounds().Dy())
}
