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
	return (strconv.FormatUint(n, 36))
}

func createImage(formContent string) string {
	fontPath := "src/fonts/Koruri.ttf"
	bgimagePath := "src/image/ema_bg.png"
	// bgColor := color.RGBA{ 255, 255, 255, 1 }
	// textColor := color.RGBA{ 0, 0, 0, 1 }

	// hexColor := "fff"
	// c, err := text2img.Hex(hexColor)
	// fmt.Print(c)

  d, err := text2img.NewDrawer(text2img.Params{
		FontPath: fontPath,
		BackgroundImagePath: bgimagePath,
		// BackgroundColor: c,
		// BackgroundColor: bgColor,
		// TextColor: textColor,
	})
	fmt.Print(err)

  img, err := d.Draw(formContent)
  fmt.Print(err)

	fileName := randomString() + ".jpg"
  file, err := os.Create("dist/" + fileName)
  fmt.Print(err)
  defer file.Close()

  err = jpeg.Encode(file, img, &jpeg.Options{Quality: 70})
  fmt.Print(err)
	fmt.Print(fileName)
	return (fileName)
} 

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")

	r.POST("/", func(c *gin.Context) {
		formContent := c.PostForm("content")
		fileName := createImage(formContent)
		c.HTML(http.StatusOK, "result.tmpl", gin.H{
			"content": formContent,
			"fileName": fileName,
		})
	})

	r.Run(":8080")
}
