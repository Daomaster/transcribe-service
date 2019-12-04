package api

import (
	"github.com/Daomaster/transcribe-service/api/middleware"
	"github.com/Daomaster/transcribe-service/api/transcription"
	"github.com/Daomaster/transcribe-service/api/user"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// get the auth middleware
	authMiddleware, err := middleware.GetAuthMiddleware()
	if err != nil {
		logrus.Fatal(err)
	}

	// user routes
	userRoute := r.Group(`/api/users`)
	userRoute.Use()
	{
		// create user
		userRoute.POST("", user.CreateUser)
	}

	// auth routes
	authRoute := r.Group(`/api/auth`)
	authRoute.Use()
	{
		// login
		authRoute.POST(`/login`, authMiddleware.LoginHandler)
		// refresh token
		authRoute.GET(`/refresh`, authMiddleware.RefreshHandler)
	}

	// transcription routes
	transcriptionRoute := r.Group(`/api/transcription`)
	transcriptionRoute.Use(authMiddleware.MiddlewareFunc())
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
