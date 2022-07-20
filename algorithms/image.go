package algorithms

import (
	"image"
	"image/color"
)

func AverageColor(img image.Image) [3]float64 {
	bounds := img.Bounds()
	r, g, b := 0.0, 0.0, 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r, g, b = r+float64(r1), g+float64(g1), b+float64(b1)
		}
	}
	//总像素
	totalPixels := float64(bounds.Max.X * bounds.Max.Y)
	//平均像素
	return [3]float64{r / totalPixels, g / totalPixels, b / totalPixels}
}

func zoomSize(img image.Image, newWidth int) image.NRGBA {
	bounds := img.Bounds()
	width := bounds.Dx()
	//缩放比率
	ratio := width / newWidth
	out := image.NewNRGBA(image.Rect(bounds.Min.X/ratio, bounds.Min.X/ratio, bounds.Max.X/ratio, bounds.Max.Y/ratio))
	for y, j := bounds.Min.Y, bounds.Min.Y; y < bounds.Max.Y; y, j = y+ratio, j+1 {
		for x, i := bounds.Min.X, bounds.Min.X; x < bounds.Max.X; x, i = x+ratio, i+1 {
			r, g, b, a := img.At(x, y).RGBA()
			out.SetNRGBA(i, j, color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)})
		}
	}
	return *out
}
