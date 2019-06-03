package main

import (
	"golang.org/x/tour/pic"
	"golang.org/x/tour/reader"
	"image"
	"image/color"
	"io"
	"os"
	"strings"
)

type MyReader struct{}

// TODO: Add a Read([]byte) (int, error) method to MyReader.

//	implement a Reader type that emits an infinite stream of the ASCII character 'A'.
func (r MyReader) Read(b []byte) (int, error) {
	l := len(b)
	for i := range b {
		b[i] = 'A'
	}
	return l, nil
}

func readerTest() {
	reader.Validate(MyReader{})
}

type rot13Reader struct {
	r io.Reader
}

func rot13(sb byte) byte {
	s := rune(sb)
	if s >= 'a' && s <= 'm' || s >= 'A' && s <= 'M' {
		sb += 13
	}
	if s >= 'n' && s <= 'z' || s >= 'N' && s <= 'Z' {
		sb -= 13
	}

	return sb
}

func (reader rot13Reader) Read(b []byte) (int, error) {
	n, err := reader.r.Read(b)
	for i := 0; i < n; i++ {
		b[i] = rot13(b[i])
	}
	return n, err
}

func exerciseRotReader() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}

type Image struct {
	x, y int
	data [][]color.Color
}

func (img Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (img Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.x, img.y)
}

func (img Image) At(x, y int) color.Color {
	return img.data[y][x]
}

func imageExercise() {
	m := NewImage(100, 100)
	pic.ShowImage(m)
}

func NewImage(dx, dy int) Image {
	img := Image{dx, dy, make([][]color.Color, dy)}
	for y := range img.data {
		sx := make([]color.Color, dx)
		for x := range sx {
			//value := (x + y) / 2
			value := uint8(x * y)
			//value := x ^ y
			sx[x] = color.RGBA{value, value, 255, 255}
		}
		img.data[y] = sx
	}

	return img
}
