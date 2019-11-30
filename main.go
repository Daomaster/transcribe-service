package main

import (
	"github.com/Daomaster/transcribe-service/api"
	"github.com/Daomaster/transcribe-service/config"
	"github.com/Daomaster/transcribe-service/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	config.InitConfig()
	models.InitModel()
}

func main() {
	// init router
	g := api.InitRouter()
	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	// run on 8888 for the server
	err := g.Run(":8888")
	if err != nil {
		logrus.Fatal(err)
	}
}
