package model

import (
	"database/sql"
	"fmt"
	"net/http"
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
	// TODO: 白名单校验

	var res NewVersion
	// defaultVersion := NewVersion{ // 未命中
	// 	"https://download.com",
	// 	"1.2.3.4",
	// 	"asdfghj",
	// 	"新版本",
	// 	"这是一个新版本测试信息",
	// }

	v.Update_version_code = EncodeVersion(v.Update_version_code) // 改写版本编码

	if v.Device_platform == "iOS" {
		queryStr := "SELECT id, update_version_code,download_url,md5,title,update_tips FROM rulesforios WHERE aid=? AND cpu_arch=? AND channel=? AND max_update_version_code>=? AND min_update_version_code<=? order by update_version_code DESC"
		rows, err := db.Query(queryStr, v.Aid, v.Cpu_arch, v.Channel, v.Update_version_code, v.Update_version_code)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&res.Update_version_code, &res.Download_url, &res.Md5, &res.Title, &res.Update_tips)
			if err != nil {
				fmt.Printf("failed2, err:%v\n", err)
				return nil, err
			}
			if check_ID(v, &res, db) { //在白名单上
				v.Update_version_code = DecodeVersion(v.Update_version_code)
				res.Download_url = "/download?aid=" + strconv.Itoa(v.Aid) + "&platform=" + v.Device_platform + "&update_version_code=" + res.Update_version_code + "&url=" + res.Download_url
				return &res, nil
			}
		}

	} else {
		queryStr := "SELECT id, update_version_code,download_url,md5,title,update_tips FROM rulesforandroid WHERE aid=? AND cpu_arch=? AND channel=? AND max_update_version_code>=? AND min_update_version_code<=? AND max_os_api>=? AND min_os_api<=? order by update_version_code DESC"
		rows, err := db.Query(queryStr, v.Aid, v.Cpu_arch, v.Channel, v.Update_version_code, v.Update_version_code, v.Os_api, v.Os_api)
		if err != nil {
			fmt.Printf("failed1, err:%v\n", err)
			return nil, err
		}
		defer rows.Close()

		var id int
		for rows.Next() {
			err := rows.Scan(&id, &res.Update_version_code, &res.Download_url, &res.Md5, &res.Title, &res.Update_tips)
			if err != nil {
				fmt.Printf("failed2, err:%v\n", err)
				return nil, err
			}

			id += 123 // 加密 id

			if check_ID(v, &res, db) { //在白名单上
				v.Update_version_code = DecodeVersion(v.Update_version_code)
				res.Download_url = "/download?id=" + strconv.Itoa(id) + "&platform=" + v.Device_platform + "&url=" + res.Download_url
				return &res, nil
			}
		}
	}

	v.Update_version_code = DecodeVersion(v.Update_version_code)
	return nil, nil
}

func GetNewVersion(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*") // 允许跨域

	version, err := GetVersion(c) // 获取提交的版本信息

	if err != nil { // 提交的版本信息有误
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid parameter"})
		return
	}

	db, err := ConnectDatabase()
	if err != nil { // 数据库连接失败
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "database error"})
		return
	}
	defer db.Close()

	newVersion, err := MatchRule(version, db) // 获取可更新新版本
	if err != nil {                           // 数据库查询失败
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "database error"})
		return
	}

	c.JSON(http.StatusOK, newVersion) // 将可更新新版本写入返回的消息
}
