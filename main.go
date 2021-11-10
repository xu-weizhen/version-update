package main

import (
	"net/http"
	"team10.com/version-update/model"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*") // html 文件
	// r.LoadHTMLFiles("template/homePage.html")   // html 文件
	r.Static("/static", "./static") // 静态文件映射
	// r.StaticFS("/favicon.ico", "./static/favicon.ico")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "homePage.html", nil) // 将 html 文件写入消息体
	})

	r.GET("/config", func(c *gin.Context) {
		c.HTML(http.StatusOK, "configPage.html", nil) // 将 html 文件写入消息体
	})

	r.POST("/config", model.GetAllRule)

	r.GET("/newversion", model.GetNewVersion)

	r.GET("/download", func(c *gin.Context) {
		db, err := model.ConnectDatabase()
		if err != nil { // 数据库连接失败
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "database error"})
			return
		}
		defer db.Close()

		err = model.DownloadCount(c, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid parameter"})
		} else {
			c.Redirect(http.StatusMovedPermanently, c.Query("url"))
		}
	})

	r.Run()
}
