package models

import (
	"encoding/base64"
	"errors"
	"github.com/jinzhu/gorm"
)

var (
	ErrUserAlreadyExist = errors.New("the user already exist")
)

// gorm model for user
type User struct {
	ID       int64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

// function to create an user
func CreateUser(username string, password string) (int64, error) {
	// encode the password
	pEnc := base64.StdEncoding.EncodeToString([]byte(password))

	// make sure the username does not exist already
	var check User
	if err := db.Model(User{}).Where("username = ?", username).First(&check).Error; err != gorm.ErrRecordNotFound {
		return 0, ErrUserAlreadyExist
	}

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

// function to validate username and password, return userID
// since auto increment starts from 1 so returning 0 with no error means invalid access
func ValidateUser(username string, password string) (int64, error) {
	// get user
	var user User
	err := db.Model(User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		// user does not exist
		if err == gorm.ErrRecordNotFound {
			return 0, err
		}

		return 0, err
	}

	// check the password
	pDec, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		return 0, err
	}

	// invalid password
	if password != string(pDec) {
		return 0, nil
	}

	return user.ID, nil
}
