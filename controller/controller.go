package controller

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	"image/jpeg"
	"mosaic/algorithms"
	"net/http"
	"strconv"
	"time"
)

var TilesMap map[string][3]float64

func init() {
	TilesMap = algorithms.CreateTilesMap()
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", nil)
}

func Mosaic(c *gin.Context) {
	t0 := time.Now()

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}
	file, _ := fileHeader.Open()
	tileSize, _ := strconv.Atoi(c.PostForm("tile_size"))
	original, _, _ := image.Decode(file)
	bounds := original.Bounds()

	cloneTails := algorithms.CloneTilesDB(TilesMap)

	// fan-out
	c1 := algorithms.Cut(original, &cloneTails, tileSize, bounds.Min.X, bounds.Min.Y, bounds.Max.X/2, bounds.Max.Y/2)
	c2 := algorithms.Cut(original, &cloneTails, tileSize, bounds.Max.X/2, bounds.Min.Y, bounds.Max.X, bounds.Max.Y/2)
	c3 := algorithms.Cut(original, &cloneTails, tileSize, bounds.Min.X, bounds.Max.Y/2, bounds.Max.X/2, bounds.Max.Y)
	c4 := algorithms.Cut(original, &cloneTails, tileSize, bounds.Max.X/2, bounds.Max.Y/2, bounds.Max.X, bounds.Max.Y)

	// fan-in
	full := algorithms.Combine(bounds, c1, c2, c3, c4)
	buf1 := new(bytes.Buffer)
	jpeg.Encode(buf1, original, nil)
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())

	t1 := time.Now()

	c.HTML(http.StatusOK, "result.html", gin.H{
		"original": originalStr,
		"mosaic":   <-full,
		"duration": fmt.Sprintf("%v ", t1.Sub(t0)),
	})
}
