package algorithms

import (
	"sync"
)

type Tiles struct {
	mutex    *sync.Mutex
	tilesMap map[string][3]float64
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
