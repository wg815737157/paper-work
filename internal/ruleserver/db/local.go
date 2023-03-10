package db

import (
	"github.com/wg815737157/paper-work/config/ruleconfig"
	"github.com/wg815737157/paper-work/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var localDB *gorm.DB

func InitDB() {
	var err any
	dsn := ruleconfig.GlobalConfig.DefaultDb.DSN()
	localDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.SugarLogger().Error(err)
		return
	}
	sqldb, err := localDB.DB()
	if err != nil {
		log.SugarLogger().Error(err)
		return
	}
	sqldb.SetConnMaxLifetime(3600 * time.Second)
	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(1000)
}
func GetLocalDB() *gorm.DB {
	return localDB
}
