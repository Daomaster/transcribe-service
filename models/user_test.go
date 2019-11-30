package models

import (
	"encoding/base64"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	username = "testname"
	password = "testpass"
)

// test to make sure the create user query is correct
func TestCreateUser(t *testing.T) {
	a := assert.New(t)

	InitMockModel()

	// mock the insert package
	mocket.Catcher.NewMock().WithQuery(`INSERT  INTO "users" ("username","password") VALUES (?,?)`)
	defer mocket.Catcher.Reset()

	err := CreateUser(username, password)

	// check the results
	a.Nil(err, "there should be no errors")
}

// test to make sure there will error return when exception from db
func TestCreateUser_Exception(t *testing.T) {
	a := assert.New(t)

	InitMockModel()

	// mock the insert package
	mocket.Catcher.NewMock().WithQuery(`INSERT  INTO "users" ("username","password") VALUES (?,?)`).WithExecException()
	defer mocket.Catcher.Reset()

	err := CreateUser(username, password)

	// check the results
	a.NotNil(err, "there should be error")
}

// test to make sure the validate user works when provide correct creds
func TestValidateUser_Correct(t *testing.T) {
	a := assert.New(t)

	InitMockModel()
	pEnc := base64.StdEncoding.EncodeToString([]byte(password))
	reply := []map[string]interface{}{{"username": username, "password": pEnc}}

	// mock the insert package
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "users"  WHERE (username = testname) ORDER BY "users"."id" ASC LIMIT 1`).WithReply(reply)
	defer mocket.Catcher.Reset()

	result, err := ValidateUser(username, password)

	// check the results
	a.Nil(err, "there should not be error")
	a.Equal(true, result, "validation should pass")
}

// test to make sure the validate user works when provide incorrect creds
func TestValidateUser_Incorrect(t *testing.T) {
	a := assert.New(t)

	InitMockModel()
	pEnc := base64.StdEncoding.EncodeToString([]byte("random"))
	reply := []map[string]interface{}{{"username": username, "password": pEnc}}

	// mock the insert package
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "users"  WHERE (username = testname) ORDER BY "users"."id" ASC LIMIT 1`).WithReply(reply)
	defer mocket.Catcher.Reset()

	result, err := ValidateUser(username, password)

	// check the results
	a.Nil(err, "there should not be error")
	a.Equal(false, result, "validation should not pass")
}

// test to make sure the validate return error upon exception from db
func TestValidateUser_Exception(t *testing.T) {
	a := assert.New(t)

	InitMockModel()

	// mock the insert package
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "users"  WHERE (username = testname) ORDER BY "users"."id" ASC LIMIT 1`).WithQueryException()
	defer mocket.Catcher.Reset()

	result, err := ValidateUser(username, password)

	// check the results
	a.NotNil(err, "there should be error")
	a.Equal(false, result, "validation should not pass")
}
