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

	r.GET("/newversion", func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*") // 允许跨域

		version, err := model.GetVersion(c) // 获取提交的版本信息

		if err != nil { // 提交的版本信息有误
			c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid parameter"})
			return
		}

		db, err := model.ConnectDatabase()
		if err != nil { // 数据库连接失败
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "database error"})
			return
		}
		defer db.Close()

		newVersion, err := model.MatchRule(version, db) // 获取可更新新版本
		if err != nil {                                 // 数据库查询失败
			c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid parameter"})
			return
		}

		c.JSON(http.StatusOK, newVersion) // 将可更新新版本写入返回的消息
	})

	r.POST("/config", model.GetAllRule)

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
