package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/error_tracer"
	"user-service/service"
)

func tokenVerify(c *gin.Context) {
	token := c.Query("token")
	tokenService := service.TokenService{}

	claim, err := tokenService.Decode(token)
	if err != nil {
		error_tracer.Client.ErrorLog("tokenVerify", "invalidToken", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"data":    err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Access granted",
		"data":    claim.User,
	})
}
