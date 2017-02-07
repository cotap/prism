package prism

import (
	"crypto/sha1"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResize(t *testing.T) {
	img := testImg("mlk.png")
	_ = img.Resize(100, 100)

	assert.Equal(t, "a8e50040e1a0219b3f2dd7710917f101048e071e", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 100, img.Bounds().Dx())
	assert.Equal(t, 100, img.Bounds().Dy())
}

func TestFit(t *testing.T) {
	img := testImg("mlk.png")
	_ = img.Fit(100, 100)

	assert.Equal(t, "663fc1b903e0e4090f926687cc55c6829f57fa37", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 100, img.Bounds().Dx())
	assert.Equal(t, 96, img.Bounds().Dy())
}

func TestReorient1(t *testing.T) {
	img := testImg("orientations/orientation-1.jpg")
	_ = img.Reorient()

	assert.Equal(t, "06d6a0c8dbdc3229b8ea8d1fe257c890baa41088", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 480, img.Bounds().Dx())
	assert.Equal(t, 640, img.Bounds().Dy())
}

func TestReorient2(t *testing.T) {
	img := testImg("orientations/orientation-2.jpg")
	_ = img.Reorient()

	assert.Equal(t, "0b12c9d9c59544298929fb229cb6cdcefd2e8432", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 480, img.Bounds().Dx())
	assert.Equal(t, 640, img.Bounds().Dy())
}

func TestReorient3(t *testing.T) {
	img := testImg("orientations/orientation-3.jpg")
	_ = img.Reorient()

	assert.Equal(t, "3802560db62d5703a8d96c821efd3e839995e681", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 480, img.Bounds().Dx())
	assert.Equal(t, 640, img.Bounds().Dy())
}

func TestReorient4(t *testing.T) {
	img := testImg("orientations/orientation-4.jpg")
	_ = img.Reorient()

	assert.Equal(t, "58de3d130761d7837bc61fad72673d4ef328f18f", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 480, img.Bounds().Dx())
	assert.Equal(t, 640, img.Bounds().Dy())
}

func TestReorient5(t *testing.T) {
	img := testImg("orientations/orientation-5.jpg")
	_ = img.Reorient()

	assert.Equal(t, "101f3dc5b9834dbf00db0f0f77f71233bb6defc5", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 480, img.Bounds().Dx())
	assert.Equal(t, 640, img.Bounds().Dy())
}

func TestReorient6(t *testing.T) {
	img := testImg("orientations/orientation-6.jpg")
	_ = img.Reorient()

	assert.Equal(t, "0abcaf927976b94134156e6c0385ea3262a66457", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 480, img.Bounds().Dx())
	assert.Equal(t, 640, img.Bounds().Dy())
}

func TestReorient7(t *testing.T) {
	img := testImg("orientations/orientation-7.jpg")
	_ = img.Reorient()

	assert.Equal(t, "bcf54553b6b7a0157caf25f56147256c75298706", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 480, img.Bounds().Dx())
	assert.Equal(t, 640, img.Bounds().Dy())
}

func TestReorient8(t *testing.T) {
	img := testImg("orientations/orientation-8.jpg")
	_ = img.Reorient()

	assert.Equal(t, "22e656fb927de87e9081c14d34375bb960a9444b", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 480, img.Bounds().Dx())
	assert.Equal(t, 640, img.Bounds().Dy())
}

func TestRotate90(t *testing.T) {
	img := testImg("mlk.png")
	_ = img.Rotate90()

	assert.Equal(t, "26837accfb14ffb4c40e4ae431a832325d304791", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 504, img.Bounds().Dx())
	assert.Equal(t, 525, img.Bounds().Dy())
}

func TestRotate180(t *testing.T) {
	img := testImg("mlk.png")
	_ = img.Rotate180()

	assert.Equal(t, "017080651c5e51d8d07ea49de722b3f84f4b271f", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 525, img.Bounds().Dx())
	assert.Equal(t, 504, img.Bounds().Dy())
}

func TestRotate270(t *testing.T) {
	img := testImg("mlk.png")
	_ = img.Rotate270()

	assert.Equal(t, "302d7317bafe41022c7e0bf7020b417e586c527e", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 504, img.Bounds().Dx())
	assert.Equal(t, 525, img.Bounds().Dy())
}

func TestFlipH(t *testing.T) {
	img := testImg("mlk.png")
	_ = img.FlipH()

	assert.Equal(t, "92dc081f5d70eb0d2d103419dec3d6ffb416ac44", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 525, img.Bounds().Dx())
	assert.Equal(t, 504, img.Bounds().Dy())
}

func TestFlipV(t *testing.T) {
	img := testImg("mlk.png")
	_ = img.FlipV()

	assert.Equal(t, "3d756cbb38ec4327986d84e8971093cda2712e39", fmt.Sprintf("%x", sha1.Sum(img.Bytes())))
	assert.Equal(t, 525, img.Bounds().Dx())
	assert.Equal(t, 504, img.Bounds().Dy())
}

func BenchmarkResizeDown(b *testing.B) {
	mlk := testImg("mlk.png")
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		img := mlk.Copy()
		_ = img.Resize(100, 100)
		img.Release()
	}
}

func BenchmarkResizeUp(b *testing.B) {
	mlk := testImg("mlk.png")
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		img := mlk.Copy()
		_ = img.Resize(1000, 1000)
		img.Release()
	}
}

func BenchmarkFit(b *testing.B) {
	mlk := testImg("mlk.png")
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		img := mlk.Copy()
		_ = img.Fit(100, 100)
		img.Release()
	}
}

func BenchmarkRotate90(b *testing.B) {
	mlk := testImg("mlk.png")
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		img := mlk.Copy()
		_ = img.Rotate90()
		img.Release()
	}
}

func BenchmarkRotate180(b *testing.B) {
	mlk := testImg("mlk.png")
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		img := mlk.Copy()
		_ = img.Rotate180()
		img.Release()
	}
}

func BenchmarkRotate270(b *testing.B) {
	mlk := testImg("mlk.png")
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		img := mlk.Copy()
		_ = img.Rotate270()
		img.Release()
	}
}

func BenchmarkFlipH(b *testing.B) {
	mlk := testImg("mlk.png")
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		img := mlk.Copy()
		_ = img.FlipH()
		img.Release()
	}
}

func BenchmarkFlipV(b *testing.B) {
	mlk := testImg("mlk.png")
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		img := mlk.Copy()
		_ = img.FlipH()
		img.Release()
	}
}
