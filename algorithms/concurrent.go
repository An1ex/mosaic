package algorithms

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"sync"
)

func Cut(original image.Image, tiles *Tiles, tileSize, x1, y1, x2, y2 int) <-chan image.Image {
	c := make(chan image.Image)
	sp := image.Point{X: 0, Y: 0}

	// goroutine
	go func() {
		newImage := image.NewNRGBA(image.Rect(x1, y1, x2, y2))
		for y := y1; y < y2; y = y + tileSize {
			for x := x1; x < x2; x = x + tileSize {
				r, g, b, _ := original.At(x, y).RGBA()
				rgb := [3]float64{float64(r), float64(g), float64(b)}
				nearest := tiles.nearest(rgb)
				file, err := os.Open(nearest)
				if err == nil {
					img, _, err := image.Decode(file)
					if err == nil {
						t := zoomSize(img, tileSize)
						tile := t.SubImage(t.Bounds())
						tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
						draw.Draw(newImage, tileBounds, tile, sp, draw.Src)
					} else {
						//fmt.Println("error in decoding nearest", err, nearest)
						panic(err)
					}
				} else {
					//fmt.Println("error opening file when creating mosaic:", nearest)
					panic(err)
				}
				file.Close()
			}
		}
		c <- newImage.SubImage(newImage.Rect)
	}()

	return c
}

func Combine(r image.Rectangle, c1, c2, c3, c4 <-chan image.Image) <-chan string {
	c := make(chan string)

	// goroutine
	go func() {
		var wg sync.WaitGroup
		newImage := image.NewNRGBA(r)
		splice := func(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
			draw.Draw(dst, r, src, sp, draw.Src)
			wg.Done()
		}
		wg.Add(4)
		var s1, s2, s3, s4 image.Image
		var ok1, ok2, ok3, ok4 bool
		for {
			select {
			case s1, ok1 = <-c1:
				go splice(newImage, s1.Bounds(), s1, image.Point{X: r.Min.X, Y: r.Min.Y})
			case s2, ok2 = <-c2:
				go splice(newImage, s2.Bounds(), s2, image.Point{X: r.Max.X / 2, Y: r.Min.Y})
			case s3, ok3 = <-c3:
				go splice(newImage, s3.Bounds(), s3, image.Point{X: r.Min.X, Y: r.Max.Y / 2})
			case s4, ok4 = <-c4:
				go splice(newImage, s4.Bounds(), s4, image.Point{X: r.Max.X / 2, Y: r.Max.Y / 2})
			}
			if ok1 && ok2 && ok3 && ok4 {
				break
			}
		}

		// wait till all splice goroutines are complete
		wg.Wait()
		buf2 := new(bytes.Buffer)
		jpeg.Encode(buf2, newImage, nil)
		c <- base64.StdEncoding.EncodeToString(buf2.Bytes())
	}()
	return c
}
