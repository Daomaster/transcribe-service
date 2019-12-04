package models

import (
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
)

// struct for mocking the return of the join query with the correct order
type MockTranscriptionDBModel struct {
	ID       string `json:"id"`
	FilePath string `json:"file_path"`
	Result   string `json:"result"`
	FileName string `json:"file_name"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

// change the gorm client to use the mock db for testing
func InitMockModel() {
	mocket.Catcher.Register()
	db, _ = gorm.Open(mocket.DriverName, "connection_string")
	mocket.Catcher.Logging = true
}
