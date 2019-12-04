package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"testing"
)

// test for create transcription
func TestCreateTranscription(t *testing.T) {
	a := assert.New(t)

	InitMockModel()

	var expectedId = "id1"

	mocket.Catcher.NewMock().WithQuery(`INSERT  INTO "transcriptions" ("id","user_id","file_path","file_name","result") VALUES (?,?,?,?,?)`)
	defer mocket.Catcher.Reset()

	id, err := CreateTranscription(expectedId, "path", 1,"name", "json")

	// check the results
	a.Nil(err, "there should be no error")
	a.Equal(expectedId, id, "the id should match as expected")
}

// test for create transcription when there is exception from the database
func TestCreateTranscription_Exception(t *testing.T) {
	a := assert.New(t)

	InitMockModel()

	mocket.Catcher.NewMock().WithQuery(`INSERT  INTO "transcriptions" ("id","user_id","file_path","file_name","result") VALUES (?,?,?,?,?)`).
		WithExecException()
	defer mocket.Catcher.Reset()

	_, err := CreateTranscription("id1", "path", 1,"name", "json")

	// check the results
	a.NotNil(err, "there should be error")
}

// test for get transcription
func TestGetTranscription(t *testing.T) {
	a := assert.New(t)

	InitMockModel()

	// set up expected result
	transcription1 := MockTranscriptionDBModel{
		ID:       "id1",
		FilePath: "testpath1",
		FileName: "1.mp4",
		Result:   "json1",
		UserID:   1,
		Username: "user1",
	}

	transcription2 := MockTranscriptionDBModel{
		ID:       "id2",
		FilePath: "testpath2",
		FileName: "2.mp4",
		Result:   "json2",
		UserID:   2,
		Username: "user2",
	}
	transcription := []MockTranscriptionDBModel{transcription1, transcription2}

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(transcription)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that query the orders
	mocket.Catcher.NewMock().WithQuery(`SELECT transcriptions.id, transcriptions.file_path, transcriptions.result, transcriptions.file_name, user_id, username FROM "transcriptions" left join users on users.id = transcriptions.user_id`).
		WithReply(expectMap)
	defer mocket.Catcher.Reset()

	results, err := GetTranscription()

	// check if the results match
	a.Nil(err, "should be no error")
	a.Equal(len(transcription), len(results), "length of the result should be expected")
	a.Equal(transcription1.Result, (*results[0]).Result, "trans 1 should match")
	a.Equal(transcription2.Result, (*results[1]).Result, "trans 2 should match")
}

// test for get transcription when there is exception from database
func TestGetTranscription_Exception(t *testing.T) {
	a := assert.New(t)

	InitMockModel()

	// mock the query that query the orders
	mocket.Catcher.NewMock().WithQuery(`SELECT transcriptions.id, transcriptions.file_path, transcriptions.result, transcriptions.file_name, user_id, username FROM "transcriptions" left join users on users.id = transcriptions.user_id`).
		WithQueryException()
	defer mocket.Catcher.Reset()

	_, err := GetTranscription()

	// check if the results match
	a.NotNil(err, "should be error")
}

// test for get transcription by id
func TestGetTranscriptionByID(t *testing.T) {
	a := assert.New(t)

	InitMockModel()

	// set up expected result
	transcription1 := MockTranscriptionDBModel{
		ID:       "id1",
		FilePath: "testpath1",
		Result:   "json1",
		FileName: "1.mp4",
		UserID:   1,
		Username: "user1",
	}

	transcription := []MockTranscriptionDBModel{transcription1}

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(transcription)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that query the orders
	mocket.Catcher.NewMock().
		WithQuery(`SELECT transcriptions.id, transcriptions.file_path, transcriptions.result, transcriptions.file_name, user_id, username FROM "transcriptions" left join users on users.id = transcriptions.user_id WHERE (transcriptions.id = id1)`).
		WithReply(expectMap)
	defer mocket.Catcher.Reset()

	result, err := GetTranscriptionByID(transcription1.ID)

	// check if the results match
	a.Nil(err, "should be no error")
	a.Equal(transcription1.ID, result.ID, "result id should match")
	a.Equal(transcription1.Result, result.Result, "json result should match")
}

// test for get transcription by id
func TestGetTranscriptionByID_Not_Found(t *testing.T) {
	a := assert.New(t)

	InitMockModel()

	// set up expected result
	transcription1 := Transcription{
		ID:       "id1",
		UserID:   1,
		FilePath: "testpath1",
		Result:   "json1",
	}

	transcription := []Transcription{transcription1}

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(transcription)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that query the orders
	mocket.Catcher.NewMock().
		WithQuery(`SELECT transcriptions.id, transcriptions.file_path, transcriptions.result, transcriptions.file_name, user_id, username FROM "transcriptions" left join users on users.id = transcriptions.user_id WHERE (transcriptions.id = id1)`).
		WithError(gorm.ErrRecordNotFound)
	defer mocket.Catcher.Reset()

	_, err := GetTranscriptionByID(transcription1.ID)

	// check if the results match
	a.NotNil(err, "should be error")
	a.Equal(gorm.ErrRecordNotFound, err, "error should match")
}

// test for get transcription by id upon exception
func TestGetTranscriptionByID_Exception(t *testing.T) {
	a := assert.New(t)

	InitMockModel()

	// mock the query that query the orders
	mocket.Catcher.NewMock().
		WithQuery(`SELECT transcriptions.id, transcriptions.file_path, transcriptions.result, transcriptions.file_name, user_id, username FROM "transcriptions" left join users on users.id = transcriptions.user_id WHERE (transcriptions.id = id1)`).
		WithQueryException()
	defer mocket.Catcher.Reset()

	_, err := GetTranscriptionByID("id1")

	// check if the results match
	a.NotNil(err, "should be error")
}
