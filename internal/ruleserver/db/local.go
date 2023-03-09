package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var localDB *gorm.DB

func InitDB() {
	var err any
	dsn := "root:@tcp(localhost:3306)/paper_work?parseTime=true&loc=Local"
	localDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	sqldb, err := localDB.DB()
	if err != nil {
		log.Fatalln(err)
	}
	sqldb.SetConnMaxLifetime(3600 * time.Second)
	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(1000)

}
func GetLocalDB() *gorm.DB {
	return localDB
}
