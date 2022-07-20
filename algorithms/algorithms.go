package algorithms

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sync"
)

type Tiles struct {
	mutex    *sync.Mutex
	tilesMap map[string][3]float64
}

func averageColor(img image.Image) [3]float64 {
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

func CreateTilesMap() map[string][3]float64 {
	tilesMap := make(map[string][3]float64)
	files, _ := ioutil.ReadDir("tilesDB")
	for _, f := range files {
		name := filepath.Join("tilesDB", f.Name())
		file, err := os.Open(name)
		if err == nil {
			img, _, err := image.Decode(file)
			if err == nil {
				tilesMap[name] = averageColor(img)
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}
	return tilesMap
}

//寻找和目标图票平均颜色最相似的一张瓷砖图片
func (tiles *Tiles) nearest(target [3]float64) string {
	var filename string
	tiles.mutex.Lock()
	smallest := 1000000.0
	for k, v := range tiles.tilesMap {
		dist := distance(target, v)
		if dist < smallest {
			filename, smallest = k, dist
		}
	}
	delete(tiles.tilesMap, filename)
	tiles.mutex.Unlock()
	return filename
}

//计算两点之间的欧几里得距离
func distance(p1 [3]float64, p2 [3]float64) float64 {
	return math.Sqrt(sq(p2[0]-p1[0]) + sq(p2[1]-p1[1]) + sq(p2[2]-p1[2]))
}

func sq(n float64) float64 {
	return n * n
}

func CloneTilesDB(tilesMap map[string][3]float64) Tiles {
	cloneMap := make(map[string][3]float64)
	for k, v := range tilesMap {
		cloneMap[k] = v
	}
	tiles := Tiles{
		tilesMap: cloneMap,
		mutex:    &sync.Mutex{},
	}
	return tiles
}

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
