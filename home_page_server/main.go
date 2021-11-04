package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLFiles("template/homePage.html")
	r.Static("/static", "./static")
	// r.StaticFS("/favicon.ico", "./static/favicon.ico")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "homePage.html", gin.H{})
	})

	r.Run()
}
