package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		content := c.PostForm("content")
		c.JSON(http.StatusOK, gin.H{"content": content})
	})

	r.Run(":8080")
}
