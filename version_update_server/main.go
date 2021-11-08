package main

import (
	"fmt"
	"net/http"

	"version_update/model"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*") // 允许跨域

		version, err := model.GetVersion(c) // 获取提交的版本信息
		// fmt.Println(version)

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

	router.POST("/config", func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*") // 允许跨域

		rule, err := model.GetPostRule(c)
		fmt.Println(rule.Platform)
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

		err = model.AddRuleToDatabase(db, rule)
		if err != nil { // 写入数据库失败
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "success"})
	})

	router.GET("/download", func(c *gin.Context) {
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

	router.Run(":8081")
}
