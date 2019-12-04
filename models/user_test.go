package models

import (
	"encoding/base64"
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"math/rand"
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

	var expectedId = rand.Int63n(100)

	// mock the insert package
	mocket.Catcher.NewMock().WithQuery(`INSERT  INTO "users" ("username","password") VALUES (?,?)`).WithID(expectedId)
	defer mocket.Catcher.Reset()

	id, err := CreateUser(username, password)

	// check the results
	a.Nil(err, "there should be no errors")
	a.Equal(expectedId, id, "the id should match as expected")
}

// test to make sure there will error return when exception from db
func TestCreateUser_Exception(t *testing.T) {
	a := assert.New(t)

	InitMockModel()

	// mock the insert package
	mocket.Catcher.NewMock().WithQuery(`INSERT  INTO "users" ("username","password") VALUES (?,?)`).WithExecException()
	defer mocket.Catcher.Reset()

	_, err := CreateUser(username, password)

	// check the results
	a.NotNil(err, "there should be error")
}

// test to make sure the validate user works when provide correct creds
func TestValidateUser_Correct(t *testing.T) {
	a := assert.New(t)

	const expectedID = 5

	InitMockModel()
	pEnc := base64.StdEncoding.EncodeToString([]byte(password))
	reply := []map[string]interface{}{{"id": expectedID, "username": username, "password": pEnc}}

	// mock the insert package
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "users"  WHERE (username = testname) ORDER BY "users"."id" ASC LIMIT 1`).WithReply(reply)
	defer mocket.Catcher.Reset()

	result, err := ValidateUser(username, password)

	// check the results
	a.Nil(err, "there should not be error")
	a.Equal(int64(expectedID), result, "validation should pass")
}

// test to make sure the validate user works when provide incorrect creds
func TestValidateUser_Incorrect(t *testing.T) {
	a := assert.New(t)

	InitMockModel()
	pEnc := base64.StdEncoding.EncodeToString([]byte("random"))
	reply := []map[string]interface{}{{"id": 5,"username": username, "password": pEnc}}

	// mock the insert package
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "users"  WHERE (username = testname) ORDER BY "users"."id" ASC LIMIT 1`).WithReply(reply)
	defer mocket.Catcher.Reset()

	result, err := ValidateUser(username, password)

	// check the results
	a.Nil(err, "there should not be error")
	a.Equal(int64(0), result, "validation should not pass")
}

// test to make sure the validate user works when no user found
func TestValidateUser_Not_Found(t *testing.T) {
	a := assert.New(t)

	InitMockModel()

	// mock the insert package
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "users"  WHERE (username = testname) ORDER BY "users"."id" ASC LIMIT 1`).WithError(gorm.ErrRecordNotFound)
	defer mocket.Catcher.Reset()

	result, err := ValidateUser(username, password)

	// check the results
	a.NotNil(err, "there should be error")
	a.Equal(gorm.ErrRecordNotFound, err, "error should match")
	a.Equal(int64(0), result, "validation should not pass")
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
	a.Equal(int64(0), result, "validation should not pass")
}
