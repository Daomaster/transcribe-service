package models

// gorm model for user
type User struct {
	ID       int64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	UserName string
	Password string
}

// function to create an user
func CreateUser(username string, password string) error {
	return nil
}

// function to validate username and password
func ValidateUser(username string, password string) (bool, error) {
	return false, nil
}
