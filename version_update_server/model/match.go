package model

import (
	"database/sql"
	"fmt"
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

func MatchRule(v *Version, db *sql.DB) (*NewVersion, error) {
	// 根据当前版本匹配新版本信息
	// TODO

	var res NewVersion
	// defaultVersion := NewVersion{ // 未命中
	// 	"https://download.com",
	// 	"1.2.3.4",
	// 	"asdfghj",
	// 	"新版本",
	// 	"这是一个新版本测试信息",
	// }

	if v.Device_platform == "iOS" {
		queryStr := "SELECT update_version_code,download_url,md5,title,update_tips FROM rulesforios WHERE aid=? AND cpu_arch=? AND channel=? AND max_update_version_code>=? AND min_update_version_code<=? order by update_version_code DESC limit 0,1"
		err := db.QueryRow(queryStr, v.Aid, v.Cpu_arch, v.Channel, v.Update_version_code, v.Update_version_code).Scan(&res.Update_version_code, &res.Download_url, &res.Md5, &res.Title, &res.Update_tips)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return nil, err
		}
	} else {
		queryStr := "SELECT update_version_code,download_url,md5,title,update_tips FROM rulesforandroid WHERE aid=? AND cpu_arch=? AND channel=? AND max_update_version_code>=? AND min_update_version_code<=? AND max_os_api>=? AND min_os_api<=? order by update_version_code DESC limit 0,1"
		err := db.QueryRow(queryStr, v.Aid, v.Cpu_arch, v.Channel, v.Update_version_code, v.Update_version_code, v.Os_api, v.Os_api).Scan(&res.Update_version_code, &res.Download_url, &res.Md5, &res.Title, &res.Update_tips)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return nil, err
		}
	}

	res.Download_url = "/download?aid=" + strconv.Itoa(v.Aid) + "&platform=" + v.Device_platform + "&update_version_code=" + res.Update_version_code + "&url=" + res.Download_url
	return &res, nil
}
