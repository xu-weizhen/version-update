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
	// 从消息中提取提交的版本信息

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

	// os_api 字段由 string 类型转为 int 类型
	os, err := strconv.Atoi(c.Query("os_api"))
	if err != nil {
		return &v, err
	}
	v.Os_api = os

	// aid 字段由 string 类型转为 int 类型
	aid, err := strconv.Atoi(c.Query("aid"))
	if err != nil {
		return &v, err
	}
	v.Aid = aid

	// cpu_arch 字段由 string 类型转为 int 类型
	cpu, err := strconv.Atoi(c.Query("cpu_arch"))
	if err != nil {
		return &v, err
	}
	v.Cpu_arch = cpu

	return &v, nil
}

func MatchRule(v *Version) *NewVersion {
	// 根据当前版本匹配新版本信息
	// TODO

	// 虚假的新版本信息
	newVersion := NewVersion{
		"https://download.com",
		"1.2.3.4",
		"asdfghj",
		"新版本",
		"这是一个新版本测试信息",
	}

	return &newVersion
}
