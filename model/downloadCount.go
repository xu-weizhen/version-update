package model

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func DownloadCount(c *gin.Context, db *sql.DB) error {
	tx, _ := db.Begin()
	_, err := tx.Exec("UPdate downloadCount set count=count+1 where aid=? and platform=? and update_version_code=?", c.Query("aid"), c.Query("platform"), c.Query("update_version_code"))

	if err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
		return nil
	}
}
