package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLFiles("template/homePage.html") // html 文件
	r.Static("/static", "./static")           // 静态文件映射
	// r.StaticFS("/favicon.ico", "./static/favicon.ico")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "homePage.html", gin.H{}) // 将 html 文件写入消息体
	})

	r.Run()
}
