package api

import (
	"bytes"
	"encoding/json"
	"github.com/Daomaster/transcribe-service/api/user"
	"github.com/Daomaster/transcribe-service/models"
	"github.com/Daomaster/transcribe-service/pkgs/e"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	username = "testname"
	password = "testpass"
)

// helper function to create json io
func createJson(i interface{}) (io.Reader, error) {
	body, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(body), err
}

// helper function to parse json
func parseJson(b io.Reader, i interface{}) error {
	body, err := ioutil.ReadAll(b)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, i)
	if err != nil {
		return err
	}

	return nil
}

// test success from create user
func TestCreateUser(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// get the router
	r := InitRouter()

	// create request body
	var createUser user.CreateUserRequest
	createUser.Username = username
	createUser.Password = password
	reqBody, err := createJson(createUser)

	a.Nil(err, "should not have problem with create json")

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", reqBody)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusCreated, w.Code, "server should return back 201 Created")
}

// test bad request from create user
func TestCreateUser_BadRequest(t *testing.T) {
	a := assert.New(t)

	// init the mock database
	models.InitMockModel()

	// get the router
	r := InitRouter()

	// create request body
	var createUser user.CreateUserRequest
	createUser.Username = username
	reqBody, err := createJson(createUser)

	a.Nil(err, "should not have problem with create json")

	// make request to recorder
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", reqBody)
	r.ServeHTTP(w, req)

	// check response code
	a.Equal(http.StatusBadRequest, w.Code, "server should return back 400 Bad Request")

	// parsing the error response
	var errorResponse e.ResponseError
	errRes := parseJson(w.Body, &errorResponse)
	a.Nil(errRes, "should not error out upon parsing error")
	a.NotNil(errorResponse, "server should have error response")
	a.Equal(user.ErrUserRequestInvalid.Error(), errorResponse.Error, "error response should match the error content")
}