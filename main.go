package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"image/jpeg"
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

func uploadToS3(fileName string) {

	AWS_ACCESS_KEY := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucket := os.Getenv("S3_BUCKET_NAME")
	region := "us-west-2"

	sess := session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(AWS_ACCESS_KEY, AWS_SECRET_ACCESS_KEY, ""),
			Region: aws.String(region),
	}))

	uploader := s3manager.NewUploader(sess)

	f, err  := os.Open(fileName)
	if err != nil {
			fmt.Print("failed to open file")
	}

	res, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fileName),
			Body:   f,
	})

	if err != nil {
			fmt.Println(res)
			fmt.Print(err)
	}

	fmt.Print("successfully uploaded file")
	_ = os.Remove(fileName)
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

  img, err := d.Draw(formContent)

	fileName := randomString() + ".jpg"
  file, err := os.Create(fileName)
  defer file.Close()

	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 70})
	uploadToS3(fileName)
	fmt.Print("Uploaded Image")

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
