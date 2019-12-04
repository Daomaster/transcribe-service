package models

import (
	"database/sql"
	"errors"
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
	if err != nil {
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

	// get all the columns
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// scan each row in to the struct
	for rows.Next() {
		var t Transcription

		cols := make([]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			cols[i], err = columnMapper(columns[i], &t)
			if err != nil {
				return nil, err
			}
		}

		err = rows.Scan(cols...)
		if err != nil {
			return nil, err
		}

		trans = append(trans, &t)
	}

	// if no rows returned
	if len(trans) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return trans, nil
}

// helper function to map the database field to the struct properly
func columnMapper(colName string, t *Transcription) (interface{}, error) {
	// init the user property
	var user User
	t.User = &user

	switch colName {
	case "id":
		return &t.ID, nil
	case "file_path":
		return &t.FilePath, nil
	case "file_name":
		return &t.FilePath, nil
	case "result":
		return &t.Result, nil
	case "user_id":
		return &t.User.ID, nil
	case "username":
		return &t.User.Username, nil
	default:
		return nil, errors.New("unknown column")
	}
}
