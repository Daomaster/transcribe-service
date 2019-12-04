package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

// gorm model for transcription
type Transcription struct {
	ID       string `gorm:"PRIMARY_KEY" json:"id"`
	User     *User  `gorm:"FOREIGNKEY:id" json:"user"`
	UserID   int64  `json:"-"`
	FilePath string `json:"filePath"`
	FileName string `json:"fileName"`
	Result   string `json:"result" sql:"type:longtext"`
}

// function to create an transcription
func CreateTranscription(id string, storagePath string, userID int64, fileName string, result string) (string, error) {
	// create transcription
	t := Transcription{
		ID:       id,
		UserID:   userID,
		FilePath: storagePath,
		Result:   result,
		FileName: fileName,
	}

	if err := db.Model(Transcription{}).Create(&t).Error; err != nil {
		return "", err
	}

	return t.ID, nil
}

// function to get all transcription
func GetTranscription() ([]*Transcription, error) {
	// get the raw sql raws from gorm
	rows, err := db.
		Model(Transcription{}).
		Joins("left join users on users.id = transcriptions.user_id").
		Select("transcriptions.id, transcriptions.file_path, transcriptions.result, transcriptions.file_name, user_id, username").
		Rows()
	if err != nil {
		return nil, err
	}

	// get the data from the raw sql rows
	trans, err := mapTranscriptionFromRow(rows)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return trans, nil
}

// function to get transcription by id
func GetTranscriptionByID(id string) (*Transcription, error) {
	// get the raw sql raws from gorm
	rows, err := db.
		Model(Transcription{}).
		Joins("left join users on users.id = transcriptions.user_id").
		Select("transcriptions.id, transcriptions.file_path, transcriptions.result, transcriptions.file_name, user_id, username").
		Where("transcriptions.id = ?", id).
		Rows()
	if err != nil {
		return nil, err
	}

	// get the data from the raw sql rows
	trans, err := mapTranscriptionFromRow(rows)
	if err != nil {
		return nil, err
	}

	return trans[0], nil
}

// helper function to map joint table to the struct from raw sql rows
func mapTranscriptionFromRow(rows *sql.Rows) ([]*Transcription, error) {
	var trans []*Transcription

	defer rows.Close()

	// scan each row
	for rows.Next() {
		var u User
		var t Transcription
		err := rows.Scan(&t.ID, &t.FilePath, &t.Result, &t.FileName, &u.ID, &u.Username)
		if err != nil {
			return nil, err
		}

		t.User = &u
		trans = append(trans, &t)
	}

	// if no rows returned
	if len(trans) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return trans, nil
}
