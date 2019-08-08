package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"password-service/db"
	"password-service/error_tracer"
	"password-service/model"
	"strconv"
)

func addPassword(c *gin.Context) {
	var password model.Password
	if err := c.BindJSON(&password); err != nil {
		error_tracer.Client.InfoLog("addPassword", "requestBody", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"data":    err.Error(),
		})

		return
	}

	if !password.IsValid() {
		error_tracer.Client.InfoLog("addPassword", "requestBody", "Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Required fields can't be empty",
		})
		return
	}

	userId, err := strconv.ParseInt(c.GetHeader("user_id"), 10, 64)
	if err != nil || userId == 0 {
		error_tracer.Client.ErrorLog("addPassword", "userId", fmt.Sprintf("can't get user_id = %d Error : %v", userId, err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Can't map user",
			"data":    fmt.Sprintf("can't get user_id = %d Error : %v", userId, err),
		})

		return
	}

	password.UserId = int(userId)
	if err := db.Client.Save(&password); err != nil {
		error_tracer.Client.ErrorLog("addPassword", "database", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database error",
			"data":    err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successful",
	})
}

func passwordList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetHeader("user_id"), 10, 64)
	passwords, err := db.Client.GetByUserId(uint(userId))
	if err != nil {
		error_tracer.Client.ErrorLog("passwordList", "database", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "No password found",
			"data":    err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successful",
		"data":    passwords,
	})
}
