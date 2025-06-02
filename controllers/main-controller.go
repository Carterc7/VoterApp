package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ShowHomePage handles the root route and serves an HTML template
// Go only exports func's that start with uppercase
func ShowHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Home",
	})
}
