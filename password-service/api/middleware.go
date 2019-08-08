package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"password-service/error_tracer"
	"password-service/service"
)

func authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		userId, err := service.AuthService.GetUserId(token)
		if err != nil {
			error_tracer.Client.InfoLog("middleware", "token", fmt.Sprintf("%s is invalid token", token))
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized user",
			})

			c.Abort()
		}

		c.Request.Header.Add("user_id", fmt.Sprintf("%d", userId))

		c.Next()
	}
}
