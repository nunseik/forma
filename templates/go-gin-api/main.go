// Author: {{ .Author }}
// Created: {{ .Timestamp }}

package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Define a simple GET endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Define another example endpoint
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from the {{ .ProjectName }} API!",
		})
	})

	// Start the server on port 8080
	router.Run(":8080")
}
