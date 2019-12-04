package transcription

import (
	"errors"
	"github.com/Daomaster/transcribe-service/api/middleware"
	"github.com/Daomaster/transcribe-service/models"
	"github.com/Daomaster/transcribe-service/pkgs/e"
	"github.com/Daomaster/transcribe-service/services/storage"
	"github.com/Daomaster/transcribe-service/services/transcribe"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"net/http"
)

var (
	ErrTranscriptionRequestInvalid = errors.New("the transcription request is invalid")
)

// handler for create transcription
func CreateTranscription(c *gin.Context) {
	var req CreateTranscriptionRequest

	// get the user id from claim
	claims := jwt.ExtractClaims(c)
	userId := claims[middleware.UserIdKey]

	// if user id not embedded
	if userId == nil {
		c.Status(http.StatusForbidden)
	}

	// bind the form value from the request
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.CreateErr(ErrTranscriptionRequestInvalid))
		return
	}

	// get the file header
	fileHeader, err := c.FormFile("file")
	if err != nil {
		// other exceptions
		c.JSON(http.StatusInternalServerError, e.InternalError(err))
		return
	}

	// get the io.reader from the file header
	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.InternalError(err))
		return
	}

	// generate a uuid for this transcription
	id := uuid.New().String()

	// upload to s3
	storagePath, err := storage.Client.Upload(id, fileHeader.Filename, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.InternalError(err))
		return
	}

	// transcribe the video file that is stored in s3
	result, err := transcribe.Client.Transcribe(id, storagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.InternalError(err))
		return
	}

	// save to the database
	tid, err := models.CreateTranscription(id, storagePath, int64(userId.(float64)), fileHeader.Filename, result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.InternalError(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": tid,
	})
}

// handler for get transcription
func GetTranscription(c *gin.Context) {
	// get the transcription from db
	t, err := models.GetTranscription()
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.InternalError(err))
		return
	}

	c.JSON(http.StatusOK, t)
}

// handler for get transcription by ids
func GetTranscriptionByID(c *gin.Context) {
	var request GetTranscriptionByIDRequest

	// bind the uri from gin
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, e.CreateErr(ErrTranscriptionRequestInvalid))
		return
	}

	// get the transcription from db
	t, err := models.GetTranscriptionByID(request.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusInternalServerError, e.InternalError(err))
		return
	}

	c.JSON(http.StatusOK, t)
}
