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

	// Root endpoint with project metadata
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"project": "{{ .ProjectName }}",
			"author": "{{ .Author }}",
			"created": "{{ .Timestamp }}",
			"message": "Welcome to the {{ .ProjectName }} API!",
		})
	})

	// Health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

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
			"author": "{{ .Author }}",
		})
	})

	// Start the server on port 8080
	router.Run(":8080")
}
