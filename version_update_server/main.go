package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Version struct {
	Version             string
	Device_platform     string
	Device_id           string
	Os_api              string
	Channel             string
	Version_code        string
	Update_version_code string
	Aid                 string
	Cpu_arch            string
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {

		version := c.Query("version")
		device_platform := c.Query("device_platform")
		device_id := c.Query("device_id")
		os_api := c.Query("os_api")
		channel := c.Query("channel")
		version_code := c.Query("version_code")
		update_version_code := c.Query("update_version_code")
		aid := c.Query("aid")
		cpu_arch := c.Query("cpu_arch")

		msg := Version{version, device_platform, device_id, os_api, channel, version_code, update_version_code, aid, cpu_arch}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusOK, msg)
	})

	router.Run(":8081")
}
