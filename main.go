package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")
}
