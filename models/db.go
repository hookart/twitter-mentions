package models

import (
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDBConnection() *gorm.DB {
	if db == nil {

		dsn := viper.GetString("postgres") //"host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		// log.Println(dsn)
		dbLoc, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			db = dbLoc
		}
	}
	return db
}
