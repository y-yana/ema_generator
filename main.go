package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/julianshen/text2img"
)

func randomString() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}

func createImage(formContent string) {
	fontPath := "src/fonts/Koruri.ttf"
	bgimagePath := "src/image/ema_bg.png"
	// bgColor := color.RGBA{ 255, 255, 255, 1 }
	// textColor := color.RGBA{ 0, 0, 0, 1 }

  d, err := text2img.NewDrawer(text2img.Params{
		FontPath: fontPath,
		BackgroundImagePath: bgimagePath,
		// BackgroundColor: bgColor,
		// TextColor: textColor,
	})
	fmt.Print(err)

  img, err := d.Draw(formContent)
  fmt.Print(err)

  file, err := os.Create("dist/"+ randomString() + ".jpg")
  fmt.Print(err)
  defer file.Close()

  err = jpeg.Encode(file, img, &jpeg.Options{Quality: 70})
  fmt.Print(err)

} 

func main() {
	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		formContent := c.PostForm("content")
		c.JSON(http.StatusOK, gin.H{"content": formContent})
		createImage(formContent)
	})

	r.Run(":8080")
}
