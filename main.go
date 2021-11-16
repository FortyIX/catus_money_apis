package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/hi", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "okay",
		})
	})

	err := router.Run(":3990")
	if err != nil {
		return
	}
}
