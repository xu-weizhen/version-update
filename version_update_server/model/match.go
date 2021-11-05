package model

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Version struct {
	Version             string
	Device_platform     string
	Device_id           string
	Os_api              int
	Channel             string
	Version_code        string
	Update_version_code string
	Aid                 int
	Cpu_arch            int
}

type NewVersion struct {
	Download_url        string
	Update_version_code string
	Md5                 string
	Title               string
	Update_tips         string
}

func GetVersion(c *gin.Context) (*Version, error) {

	v := Version{
		c.Query("version"),
		c.Query("device_platform"),
		c.Query("device_id"),
		-1, // Os_api
		c.Query("channel"),
		c.Query("version_code"),
		c.Query("update_version_code"),
		-1, // Aid
		-1, // Cpu_arch
	}

	os, err := strconv.Atoi(c.Query("os_api"))
	if err != nil {
		return &v, err
	}
	v.Os_api = os

	aid, err := strconv.Atoi(c.Query("aid"))
	if err != nil {
		return &v, err
	}
	v.Aid = aid

	cpu, err := strconv.Atoi(c.Query("cpu_arch"))
	if err != nil {
		return &v, err
	}
	v.Cpu_arch = cpu

	return &v, nil
}

func MatchRule(v *Version) *NewVersion {
	// todo

	newVersion := NewVersion{
		"https://download.com",
		"1.2.3.4",
		"asdfghj",
		"新版本",
		"这是一个新版本测试信息",
	}

	return &newVersion
}
