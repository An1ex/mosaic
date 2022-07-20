package tiles

import (
	"image"
	"io/ioutil"
	"mosaic/algorithms"
	"os"
	"path/filepath"
)

var OriginTilesMap map[string][3]float64

func Init() {
	OriginTilesMap = CreateTilesMap()
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
				tilesMap[name] = algorithms.AverageColor(img)
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
