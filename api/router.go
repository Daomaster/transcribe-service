package api

import (
	"github.com/Daomaster/transcribe-service/api/login"
	"github.com/Daomaster/transcribe-service/api/transcription"
	"github.com/Daomaster/transcribe-service/api/user"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// user routes
	userRoute := r.Group(`/users`)

	userRoute.Use()
	{
		// create user
		userRoute.POST("", user.CreateUser)
	}

	// login routes
	loginRoute := r.Group(`/login`)

	loginRoute.Use()
	{
		// login user
		loginRoute.POST("", login.Login)
	}

	// transcription routes
	transcriptionRoute := r.Group(`/transcription`)

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
