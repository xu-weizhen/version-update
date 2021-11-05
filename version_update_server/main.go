package main

import (
	"net/http"

	"version_update/model"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")

		version, err := model.GetVersion(c)

		if err == nil {
			c.JSON(http.StatusOK, gin.H{"msg": "invalid parameter"})
		}

		newVersion := model.MatchRule(version)
		c.JSON(http.StatusOK, newVersion)
	})

	router.Run(":8081")
}
