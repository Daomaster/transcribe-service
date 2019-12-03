package user

import (
	"errors"
	"github.com/Daomaster/transcribe-service/models"
	"github.com/Daomaster/transcribe-service/pkgs/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrUserRequestInvalid = errors.New("the user request is invalid")
)

// handler for create user
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.CreateErr(ErrUserRequestInvalid))
		return
	}

	_, err := models.CreateUser(req.Username, req.Password)
	if err != nil {
		// check if user already registered
		if err == models.ErrUserAlreadyExist {
			c.JSON(http.StatusBadRequest, e.CreateErr(err))
			return
		}

		// other exceptions
		c.JSON(http.StatusInternalServerError, e.InternalError(err))
		return
	}

	c.Status(http.StatusCreated)
}
