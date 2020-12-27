package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/julianshen/text2img"
)

func randomString() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return (strconv.FormatUint(n, 36))
}

func uploadImage(fileName string) {
	AWS_ACCESS_KEY := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucket := os.Getenv("S3_BUCKET_NAME")

	creds := credentials.NewStaticCredentials(AWS_ACCESS_KEY, AWS_SECRET_ACCESS_KEY, "")
	sess, err := session.NewSession(&aws.Config{
    Credentials: creds,
    Region: aws.String("us-west-2")},
	)
	file, err := os.Open(fileName)

	defer file.Close()

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
    Bucket: aws.String(bucket),
    Key: aws.String(fileName),
    Body: file,
	})
	log.Print(err)
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

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	r.POST("/result", func(c *gin.Context) {
			formContent := c.PostForm("content")
			fileName := createImage(formContent)
		c.HTML(http.StatusOK, "result.tmpl", gin.H{
			"content": formContent,
			"fileName": fileName,
		})
	})

	port := os.Getenv("PORT")
	r.Run(":" + port)
}
