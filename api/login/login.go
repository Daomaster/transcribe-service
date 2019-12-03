package login

import (
	"github.com/Daomaster/transcribe-service/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// handler for login
func Login(c *gin.Context) (interface{}, error) {
	var loginReq UserLoginRequest
	if err := c.ShouldBind(&loginReq); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	userId, err := models.ValidateUser(loginReq.Username, loginReq.Password)
	if err != nil {
		return nil, err
	}

	if userId == 0 {
		return nil, jwt.ErrFailedAuthentication
	}

	return userId, nil
}
