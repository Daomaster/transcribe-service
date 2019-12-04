package api

import (
	"encoding/json"
	"github.com/Daomaster/transcribe-service/api/transcription"
	"github.com/Daomaster/transcribe-service/models"
	"github.com/Daomaster/transcribe-service/pkgs/e"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initTransRoute() *gin.Engine {
	r := gin.New()
	transcriptionRoute := r.Group(`/api/transcription`)
	transcriptionRoute.Use()
	{
		// create transcription
		transcriptionRoute.POST("", transcription.CreateTranscription)
		// get transcription
		transcriptionRoute.GET("", transcription.GetTranscription)
		// get specific transcription
		transcriptionRoute.GET("/:id", transcription.GetTranscriptionByID)
	}

	return r
}

// test success from get all transcription
func TestGetTranscription(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// set up expected result
	trans1 := models.MockTranscriptionDBModel{
		ID:       "id1",
		FilePath: "testpath1",
		FileName: "1.mp4",
		Result:   "json1",
		UserID:   1,
		Username: "user1",
	}

	trans2 := models.MockTranscriptionDBModel{
		ID:       "id2",
		FilePath: "testpath2",
		FileName: "2.mp4",
		Result:   "json2",
		UserID:   2,
		Username: "user2",
	}
	trans := []models.MockTranscriptionDBModel{trans1, trans2}

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(trans)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that query the orders
	mocket.Catcher.NewMock().WithQuery(`SELECT transcriptions.id, transcriptions.file_path, transcriptions.result, transcriptions.file_name, user_id, username FROM "transcriptions" left join users on users.id = transcriptions.user_id`).
		WithReply(expectMap)
	defer mocket.Catcher.Reset()

	// get the router
	r := initTransRoute()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/transcription", nil)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusOK, w.Code, "server should return back 200 OK")

	// check the response
	var transResponse []models.Transcription
	err := parseJson(w.Body, &transResponse)
	a.Nil(err, "should not error out upon parsing response")
	a.NotNil(transResponse, "server should have a response")
	a.Equal(len(trans), len(transResponse), "response should return 2 transcription")
	a.Equal(trans1.Result, transResponse[0].Result, "trans 1 should match")
	a.Equal(trans2.Result, transResponse[1].Result, "trans 2 should match")
}

// test get 500 from get all transcription
func TestGetTranscription_Internal_Error(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// mock the query that query the orders
	mocket.Catcher.NewMock().WithQuery(`SELECT transcriptions.id, transcriptions.file_path, transcriptions.result, transcriptions.file_name, user_id, username FROM "transcriptions" left join users on users.id = transcriptions.user_id`).WithQueryException()
	defer mocket.Catcher.Reset()

	// get the router
	r := initTransRoute()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/transcription", nil)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusInternalServerError, w.Code, "server should return back 500 Internal Sever Error")

	// parsing the error response
	var errorResponse e.ResponseError
	err := parseJson(w.Body, &errorResponse)
	a.Nil(err, "should not error out upon parsing error")
	a.NotNil(errorResponse, "server should have error response")
	a.Equal(e.ErrInternalError.Error(), errorResponse.Error, "error response should match the error content")
}

// test success from get transcription by id
func TestGetTranscriptionByID(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// set up expected result
	trans1 := models.MockTranscriptionDBModel{
		ID:       "id1",
		FilePath: "testpath1",
		FileName: "1.mp4",
		Result:   "json1",
		UserID:   1,
		Username: "user1",
	}

	trans := []models.MockTranscriptionDBModel{trans1}

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(trans)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that query the orders
	mocket.Catcher.NewMock().
		WithQuery(`SELECT transcriptions.id, transcriptions.file_path, transcriptions.result, transcriptions.file_name, user_id, username FROM "transcriptions" left join users on users.id = transcriptions.user_id WHERE (transcriptions.id = id1)`).
		WithReply(expectMap)
	defer mocket.Catcher.Reset()

	// get the router
	r := initTransRoute()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/transcription/id1", nil)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusOK, w.Code, "server should return back 200 OK")

	// check the response
	var transResponse models.Transcription
	err := parseJson(w.Body, &transResponse)
	a.Nil(err, "should not error out upon parsing response")
	a.NotNil(transResponse, "server should have a response")
	a.Equal(trans1.Result, transResponse.Result, "trans 1 should match")
}

// test not found from get transcription by id
func TestGetTranscriptionByID_Not_Found(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// set up expected result
	trans1 := models.MockTranscriptionDBModel{
		ID:       "id1",
		FilePath: "testpath1",
		FileName: "1.mp4",
		Result:   "json1",
		UserID:   1,
		Username: "user1",
	}

	trans := []models.MockTranscriptionDBModel{trans1}

	// make the struct into map for the database mock
	var expectMap []map[string]interface{}
	i, _ := json.Marshal(trans)
	_ = json.Unmarshal(i, &expectMap)

	// mock the query that query the orders
	mocket.Catcher.NewMock().
		WithQuery(`SELECT transcriptions.id, transcriptions.file_path, transcriptions.result, transcriptions.file_name, user_id, username FROM "transcriptions" left join users on users.id = transcriptions.user_id WHERE (transcriptions.id = id1)`).
		WithError(gorm.ErrRecordNotFound)
	defer mocket.Catcher.Reset()

	// get the router
	r := initTransRoute()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/transcription/id1", nil)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusNotFound, w.Code, "server should return back 404 not found")
}

// test get 500 from get transcription by id
func TestGetTranscriptionByID_Internal_Error(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// mock the query that query the orders
	mocket.Catcher.NewMock().
		WithQuery(`SELECT transcriptions.id, transcriptions.file_path, transcriptions.result, transcriptions.file_name, user_id, username FROM "transcriptions" left join users on users.id = transcriptions.user_id WHERE (transcriptions.id = id1)`).
		WithQueryException()
	defer mocket.Catcher.Reset()

	// get the router
	r := initTransRoute()

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/transcription/id1", nil)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusInternalServerError, w.Code, "server should return back 500 Internal Sever Error")

	// parsing the error response
	var errorResponse e.ResponseError
	err := parseJson(w.Body, &errorResponse)
	a.Nil(err, "should not error out upon parsing error")
	a.NotNil(errorResponse, "server should have error response")
	a.Equal(e.ErrInternalError.Error(), errorResponse.Error, "error response should match the error content")
}
