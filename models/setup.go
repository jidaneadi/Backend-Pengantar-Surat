package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/pengajuan_surat?charset=utf8mb4&parseTime=True&"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
}
