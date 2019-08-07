package api

import (
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
		error_tracer.Client.ErrorLog("addPassword", "requestBody", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"data":    err.Error(),
		})

		return
	}

	if !password.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Required fields can't be empty",
		})
		return
	}

	userId, err := strconv.ParseInt(c.GetHeader("user_id"), 10, 64)
	if err != nil {
		error_tracer.Client.ErrorLog("addPassword", "userId", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Can't map user",
			"data":    err.Error(),
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

func deletePassword(c *gin.Context) {
	passwordId, err := strconv.ParseInt(c.Param("password_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Can't parse password id",
			"data":    err.Error(),
		})

		return
	}
	userId, _ := strconv.ParseInt(c.GetHeader("user_id"), 10, 64)

	password, err := db.Client.GetById(uint(passwordId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "Password not found",
			"data":   err.Error(),
		})

		return
	}

	if password.UserId != int(userId) {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Password belongs to another user",
		})

		return
	}

	if err := db.Client.Delete(password); err != nil {
		error_tracer.Client.ErrorLog("deletePassword", "database", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database can't delete password",
			"status":  err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successful",
	})
}
