package main

import (
	"net/http"

	"version_update/model"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*") // 允许跨域

		version, err := model.GetVersion(c) // 获取提交的版本信息

		if err == nil { // 提交的版本信息有误
			c.JSON(http.StatusOK, gin.H{"msg": "invalid parameter"})
		}

		newVersion := model.MatchRule(version) // 获取可更新新版本
		c.JSON(http.StatusOK, newVersion)      // 将可更新新版本写入返回的消息
	})

	router.Run(":8081")
}
