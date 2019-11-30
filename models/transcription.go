package models

// gorm model for transcription
type Transcription struct {
	ID       int64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	User     User  `gorm:"FOREIGNKEY:UserID"`
	UserID   int64
	FilePath string
	Result   string
}

// function to create an transcription
func CreateTranscription(fileName string, userID int64) (*Transcription, error) {
	return nil, nil
}

// function to get all transcription
func GetTranscription() ([]Transcription, error) {
	return nil, nil
}

// function to get transcription by id
func GetTranscriptionByID() (*Transcription, error) {
	return nil, nil
}
