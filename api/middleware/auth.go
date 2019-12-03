package middleware

import (
	"github.com/Daomaster/transcribe-service/api/login"
	"github.com/Daomaster/transcribe-service/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const userIdKey = "userId"

// function to config the gin jwt middleware and return the middleware
func GetAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					userIdKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				ID: int64(claims[userIdKey].(float64)),
			}
		},
		Authenticator: login.Login,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*models.User); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.Status(http.StatusUnauthorized)
		},
		TokenLookup: "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc: time.Now,
	})

	return authMiddleware, err
}