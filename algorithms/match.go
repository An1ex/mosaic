package algorithms

import "math"

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
