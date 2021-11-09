package model

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SendToUser struct {
	Download_url        string
	Update_version_code string
	Md5                 string
	Title               string
	Update_tips         string
}

type Rule struct { //从配置接口接收的数据
	Aid                     int    //app的表示
	Platform                string //平台
	Max_update_version_code string //可升级的最大版本号
	Min_update_version_code string //可升级的最小版本号
	Max_os_api              int    //安卓的最大版本号
	Min_os_api              int    //安卓的最小版本号
	Cpu_arch                string //CPU架构
	Channel                 string //渠道号
	ToUser                  *SendToUser
}

func EncodeVersion(v string) string {
	// 将版本号 v.xx.yyy.z 编码为 000v.00xx.0yyy.000z 格式
	arr := strings.Split(v, ".")
	ret := ""

	for i := 0; i < len(arr); i++ {
		ret += strings.Repeat("0", 4-len(arr[i])) + arr[i] + "."
	}

	return ret[:len(ret)-1]
}

func DecodeVersion(v string) string {
	// 将版本号 000v.00xx.0yyy.000zv.x 解码为 x.yyy.z 格式
	arr := strings.Split(v, ".")
	ret := ""

	for i := 0; i < len(arr); i++ {
		j := 0

		for ; j < len(arr[i]) && arr[i][j] != '0'; j++ {
		}

		if j == len(arr[i]) {
			ret += "0."
		} else {
			ret += arr[i][j:] + "."
		}
	}

	return ret[:len(ret)-1]
}

func GetPostRule(c *gin.Context) (*Rule, error) {
	// 添加规则

	user := SendToUser{
		c.Query("download_url"),
		c.Query("update_version_code"),
		c.Query("md5"),
		c.Query("title"),
		c.Query("update_tips"),
	}

	rule := Rule{
		-1, // aid
		c.Query("platform"),
		c.Query("max_update_version_code"),
		c.Query("min_update_version_code"),
		-1, // max os api
		-1, // min os api
		c.Query("cpu_arch"),
		c.Query("channel"),
		&user,
	}

	// aid 字段由 string 类型转为 int 类型
	aid, err := strconv.Atoi(c.Query("aid"))
	if err != nil {
		return &rule, err
	}
	rule.Aid = aid

	// max_os_api 字段由 string 类型转为 int 类型
	max_os_api, err := strconv.Atoi(c.Query("max_os_api"))
	if err != nil {
		return &rule, err
	}
	rule.Max_os_api = max_os_api

	// min_os_api 字段由 string 类型转为 int 类型
	min_os_api, err := strconv.Atoi(c.Query("min_os_api"))
	if err != nil {
		return &rule, err
	}
	rule.Min_os_api = min_os_api

	// fmt.Println(rule.Aid)

	// err = writeDatabase(&rule)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"err": err})
	// } else {
	// 	c.JSON(http.StatusOK, "OK!")
	// }

	return &rule, nil
}

func AddRuleToDatabase(db *sql.DB, r *Rule) error {
	// 将新规则写入数据库

	// 改写版本编码
	r.Max_update_version_code = EncodeVersion(r.Max_update_version_code)
	r.Min_update_version_code = EncodeVersion(r.Min_update_version_code)
	r.ToUser.Update_version_code = EncodeVersion(r.ToUser.Update_version_code)

	err := errors.New("platform is wrong")

	// 写数据库
	if r.Platform == "iOS" {
		stmt, _ := db.Prepare(`INSERT INTO rulesForiOS(aid, platform,update_version_code,max_update_version_code,
			min_update_version_code,cpu_arch,channel,download_url,md5,title,update_tips) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
		defer stmt.Close()
		_, err = stmt.Exec(r.Aid, r.Platform, r.ToUser.Update_version_code, r.Max_update_version_code, r.Min_update_version_code,
			r.Cpu_arch, r.Channel, r.ToUser.Download_url, r.ToUser.Md5, r.ToUser.Title, r.ToUser.Update_tips)
	}

	if r.Platform == "Android" {
		stmt, _ := db.Prepare(`INSERT INTO rulesForiOS(aid, platform,update_version_code,max_update_version_code,
			min_update_version_code,cpu_arch,channel,download_url,md5,title,update_tips,max_os_api,min_os_api ) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
		defer stmt.Close()
		_, err = stmt.Exec(r.Aid, r.Platform, r.ToUser.Update_version_code, r.Max_update_version_code, r.Min_update_version_code,
			r.Cpu_arch, r.Channel, r.ToUser.Download_url, r.ToUser.Md5, r.ToUser.Title, r.ToUser.Update_tips, r.Max_os_api, r.Min_os_api)
	}

	// 还原版本编码
	r.Max_update_version_code = DecodeVersion(r.Max_update_version_code)
	r.Min_update_version_code = DecodeVersion(r.Min_update_version_code)
	r.ToUser.Update_version_code = DecodeVersion(r.ToUser.Update_version_code)

	return err
}

func ConnectDatabase() (*sql.DB, error) {
	// 连接数据库

	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	// err := db.Ping()
	return db, err
}
