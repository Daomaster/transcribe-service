package models

import (
	"encoding/base64"
	"github.com/jinzhu/gorm"
)

// gorm model for user
type User struct {
	ID       int64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Username string
	Password string
}

// function to create an user
func CreateUser(username string, password string) (int64, error) {
	// encode the password
	pEnc := base64.StdEncoding.EncodeToString([]byte(password))

	// create user
	u := User{
		Username: username,
		Password: pEnc,
	}

	if err := db.Model(User{}).Create(&u).Error; err != nil {
		return 0, err
	}

	return u.ID, nil
}

// function to validate username and password
func ValidateUser(username string, password string) (bool, error) {
	// get user
	var user User
	err := db.Model(User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		// user does not exist
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	// check the password
	pDec, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		return false, err
	}

	// invalid password
	if password != string(pDec) {
		return false, nil
	}

	return true, nil
}
