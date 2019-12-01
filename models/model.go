package models

import (
	"github.com/Daomaster/transcribe-service/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

var db *gorm.DB

// initialize the client with mysql connection
func InitModel() {
	cString := config.GetConfig().DbConfig.GetConnectionString()

	var err error
	db, err = gorm.Open("mysql", cString)
	if err != nil {
		logrus.Fatal(err)
	}
}

// initialize the tables based on the model if not exist
func migrate() {
	db.AutoMigrate(&User{}, &Transcription{})
}
