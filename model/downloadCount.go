package model

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DownloadCount(c *gin.Context) {
	db, err := ConnectDatabase()
	if err != nil { // 数据库连接失败
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "database error"})
		return
	}
	defer db.Close()

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid parameter"})
		return
	}
	id -= 123 // 解密 id

	tx, _ := db.Begin()

	// 更新下载计数
	if c.Query("platform") == "Android" {
		_, err = tx.Exec("Update rulesForAndroid set download_count=download_count+1 where id=?", id)
	} else if c.Query("platform") == "IOS" {
		_, err = tx.Exec("Update rulesForiOS set download_count=download_count+1 where id=?", id)
	} else {
		err = errors.New("invalid parameter")
	}

	if err != nil {
		tx.Rollback() // 回滚
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid parameter"})
	} else {
		tx.Commit() // 提交
		c.Redirect(http.StatusMovedPermanently, c.Query("url"))
	}
}
