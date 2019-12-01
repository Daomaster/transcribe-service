package models

import "github.com/jinzhu/gorm"

// gorm model for transcription
type Transcription struct {
	ID       int64  `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	User     *User  `gorm:"FOREIGNKEY:UserID" json:"user"`
	UserID   int64  `json:"-"`
	FilePath string `json:"file_path"`
	Result   string `json:"result"`
}

// function to create an transcription
func CreateTranscription(filepath string, userID int64, result string) (int64, error) {
	// create transcription
	t := Transcription{
		UserID:   userID,
		FilePath: filepath,
		Result:   result,
	}

	if err := db.Model(Transcription{}).Create(&t).Error; err != nil {
		return 0, err
	}

	return t.ID, nil
}

// function to get all transcription
func GetTranscription() ([]*Transcription, error) {
	var transcription []*Transcription

	err := db.Model(Transcription{}).Find(&transcription).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return transcription, nil
}

// function to get transcription by id
func GetTranscriptionByID(id int64) (*Transcription, error) {
	var transcription []*Transcription

	err := db.Model(Transcription{}).Where("id = ?", id).First(&transcription).Error
	if err != nil {
		// user does not exist
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return transcription[0], nil
}
