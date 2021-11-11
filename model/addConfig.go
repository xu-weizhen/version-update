package model

import (
	"database/sql"
	"io/ioutil"
	"net/http"
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
	Device_id_list          []byte //设备白名单
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

func GetAllRule(c *gin.Context) {
	// 添加规则

	user := SendToUser{
		c.PostForm("download_url"),
		c.PostForm("update_version_code"),
		c.PostForm("md5"),
		c.PostForm("title"),
		c.PostForm("update_tips"),
	}

	rule := Rule{
		-1, // aid
		c.PostForm("platform"),
		c.PostForm("max_update_version_code"),
		c.PostForm("min_update_version_code"),
		-1, // max os api
		-1, // min os api
		c.PostForm("cpu_arch"),
		c.PostForm("channel"),
		[]byte{}, // device_id_list
		&user,
	}

	// aid 字段由 string 类型转为 int 类型
	aid, err := strconv.Atoi(c.PostForm("aid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	rule.Aid = aid

	// max_os_api 字段由 string 类型转为 int 类型
	max_os_api, err := strconv.Atoi(c.PostForm("max_os_api"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	rule.Max_os_api = max_os_api

	// min_os_api 字段由 string 类型转为 int 类型
	min_os_api, err := strconv.Atoi(c.PostForm("min_os_api"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	rule.Min_os_api = min_os_api

	// 改写版本编码
	rule.Max_update_version_code = EncodeVersion(rule.Max_update_version_code)
	rule.Min_update_version_code = EncodeVersion(rule.Min_update_version_code)
	rule.ToUser.Update_version_code = EncodeVersion(rule.ToUser.Update_version_code)

	// 还原版本编码
	recoverVersion := func(r *Rule) {
		r.Max_update_version_code = DecodeVersion(r.Max_update_version_code)
		r.Min_update_version_code = DecodeVersion(r.Min_update_version_code)
		r.ToUser.Update_version_code = DecodeVersion(r.ToUser.Update_version_code)
	}
	defer recoverVersion(&rule) // 还原版本编码

	// 读取白名单
	f, _, err := c.Request.FormFile("device_id_list")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	} else {
		rule.Device_id_list, err = ioutil.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err})
			return
		}
		//fmt.Println(string(rule.device_id_list))
	}

	// fmt.Println(rule.Sid)

	db, err := ConnectDatabase()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Database error"})
		return
	}
	defer db.Close()

	err = AddRuleToDatabase(db, &rule)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "OK!"})

	// err = writeDatabase(&rule, &user)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
	// 	return err
	// }

	// err = writeDeviceID(&rule, &user)
	// if err != nil {
	// 	//删除已经插入数据库的数据
	// 	if rule.Platform == "iOS" {
	// 		db, err := ConnectDatabase()
	// 		if err != nil {
	// 			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
	// 			return err
	// 		}
	// 		defer db.Close()
	// 		stmt, _ := db.Prepare(`DELETE FROM rulesforiOS WHERE aid=? and platform=? and update_version_code=?`)
	// 		defer stmt.Close()
	// 		stmt.Exec(rule.Aid, rule.Platform, rule.ToUser.Update_version_code)
	// 	} else {
	// 		db, err := ConnectDatabase()
	// 		if err != nil {
	// 			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
	// 			return err
	// 		}
	// 		defer db.Close()
	// 		stmt, _ := db.Prepare(`DELETE FROM rulesforandroid WHERE aid=? and platform=? and update_version_code=?`)
	// 		defer stmt.Close()
	// 		stmt.Exec(rule.Aid, rule.Platform, rule.ToUser.Update_version_code)
	// 	}
	// 	c.JSON(http.StatusInternalServerError, gin.H{"msg": "Database error"})
	// } else {
	// 	c.JSON(http.StatusOK, "OK!")
	// }
}

// func writeDeviceID(r *Rule, u *SendToUser) error {
// 	db, err := ConnectDatabase()
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()

// 	stmt, _ := db.Prepare(`INSERT INTO device_id(aid ,platform,update_version_code,device_id_list)
// 						VALUES(?, ? , ? , ?)`)
// 	defer stmt.Close()
// 	_, err = stmt.Exec(r.Aid, r.Platform, r.ToUser.Update_version_code, string(r.Device_id_list))
// 	if err == nil {
// 		return err
// 	}
// 	return errors.New("deviceid is wrong")
// }

// func writeDatabase(r *Rule, u *SendToUser) error {
// 	db, err := ConnectDatabase()
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()

// 	if r.Platform == "iOS" {
// 		stmt, _ := db.Prepare(`INSERT INTO rulesForiOS(aid, platform,update_version_code,max_update_version_code,
// 			min_update_version_code,cpu_arch,channel,download_url,md5,title,update_tips)
// 			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
// 		defer stmt.Close()
// 		_, err := stmt.Exec(r.Aid, r.Platform, r.ToUser.Update_version_code, r.Max_update_version_code, r.Min_update_version_code,
// 			r.Cpu_arch, r.Channel, r.ToUser.Download_url, r.ToUser.Md5, r.ToUser.Title, r.ToUser.Update_tips)
// 		return err
// 	}
// 	if r.Platform == "Android" {
// 		stmt, _ := db.Prepare(`INSERT INTO rulesForiOS(aid, platform,update_version_code,max_update_version_code,
// 			min_update_version_code,cpu_arch,channel,download_url,md5,title,update_tips,max_os_api,min_os_api )
// 			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
// 		defer stmt.Close()
// 		_, err := stmt.Exec(r.Aid, r.Platform, r.ToUser.Update_version_code, r.Max_update_version_code, r.Min_update_version_code,
// 			r.Cpu_arch, r.Channel, r.ToUser.Download_url, r.ToUser.Md5, r.ToUser.Title, r.ToUser.Update_tips, r.Max_os_api, r.Min_os_api)
// 		return err
// 	}

// 	return errors.New("platform is wrong")
// }

func ConnectDatabase() (*sql.DB, error) {
	// 连接数据库

	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	// err := db.Ping()
	return db, err
}

func AddRuleToDatabase(db *sql.DB, r *Rule) error {
	// 将新规则写入数据库

	// 写数据库
	tx, err := db.Begin()

	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		return err
	}

	if r.Platform == "iOS" {
		sentence := `INSERT INTO rulesForiOS(aid, platform,update_version_code,max_update_version_code,
			min_update_version_code,cpu_arch,channel,download_url,md5,title,update_tips) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		// stmt, _ := db.Prepare(sentence)
		// defer stmt.Close()
		// _, err = stmt.Exec(r.Aid, r.Platform, r.ToUser.Update_version_code, r.Max_update_version_code, r.Min_update_version_code,
		// 	r.Cpu_arch, r.Channel, r.ToUser.Download_url, r.ToUser.Md5, r.ToUser.Title, r.ToUser.Update_tips)
		_, err = tx.Exec(sentence, r.Aid, r.Platform, r.ToUser.Update_version_code, r.Max_update_version_code, r.Min_update_version_code,
			r.Cpu_arch, r.Channel, r.ToUser.Download_url, r.ToUser.Md5, r.ToUser.Title, r.ToUser.Update_tips)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if r.Platform == "Android" {
		sentence := `INSERT INTO rulesForiOS(aid, platform,update_version_code,max_update_version_code,
			min_update_version_code,cpu_arch,channel,download_url,md5,title,update_tips,max_os_api,min_os_api ) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		// stmt, _ := db.Prepare(sentence)
		// defer stmt.Close()
		// _, err = stmt.Exec(r.Aid, r.Platform, r.ToUser.Update_version_code, r.Max_update_version_code, r.Min_update_version_code,
		// 	r.Cpu_arch, r.Channel, r.ToUser.Download_url, r.ToUser.Md5, r.ToUser.Title, r.ToUser.Update_tips, r.Max_os_api, r.Min_os_api)
		_, err = tx.Exec(sentence, r.Aid, r.Platform, r.ToUser.Update_version_code, r.Max_update_version_code, r.Min_update_version_code,
			r.Cpu_arch, r.Channel, r.ToUser.Download_url, r.ToUser.Md5, r.ToUser.Title, r.ToUser.Update_tips, r.Max_os_api, r.Min_os_api)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	sentence := `INSERT INTO device_id(aid ,platform,update_version_code,device_id_list) VALUES(?, ? , ? , ?)`
	_, err = tx.Exec(sentence, r.Aid, r.Platform, r.ToUser.Update_version_code, string(r.Device_id_list))

	if err != nil {
		if r.Platform == "iOS" {
			sentence := `DELETE FROM rulesforiOS WHERE aid=? and platform=? and update_version_code=?`
			_, err = tx.Exec(sentence, r.Aid, r.Platform, r.ToUser.Update_version_code)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			sentence := `DELETE FROM rulesforandroid WHERE aid=? and platform=? and update_version_code=?`
			_, err = tx.Exec(sentence, r.Aid, r.Platform, r.ToUser.Update_version_code)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	_, err = tx.Exec("INSERT INTO downloadcount(aid, platform, update_version_code, count) VALUES (?, ?, ?, ?)", r.Aid, r.Platform, r.ToUser.Update_version_code, 0)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
