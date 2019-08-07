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

func passwordDetails(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetHeader("user_id"), 10, 64)
	passwordId, _ := strconv.ParseInt(c.Param("password_id"), 10, 64)
	password, err := db.Client.GetById(int(passwordId), int(userId))
	if err != nil {
		error_tracer.Client.InfoLog("passwordDetails", "notFound",
			fmt.Sprintf("userId %d passwordId %d not found Error : %v", userId, passwordId, err))
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Invalid id",
			"data":    err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successful",
		"data":    password,
	})
}

func deletePassword(c *gin.Context) {
	passwordId, _ := strconv.ParseInt(c.Param("password_id"), 10, 64)
	userId, _ := strconv.ParseInt(c.GetHeader("user_id"), 10, 64)

	password, err := db.Client.GetById(int(passwordId), int(userId))
	if err != nil {
		error_tracer.Client.InfoLog("deletePassword", "notFound",
			fmt.Sprintf("userId %d passwordId %d not found Error : %v", userId, passwordId, err))
		c.JSON(http.StatusNotFound, gin.H{
			"status": "Password not found",
			"data":   err.Error(),
		})

		return
	}

	if err := db.Client.Delete(&password); err != nil {
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

func updatePassword(c *gin.Context) {
	passwordId, _ := strconv.ParseInt(c.Param("password_id"), 10, 64)
	userId, _ := strconv.ParseInt(c.GetHeader("user_id"), 10, 64)

	var password model.Password
	if err := c.BindJSON(&password); err != nil {
		error_tracer.Client.ErrorLog("updatePassword", "requestBody", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"data":    err.Error(),
		})

		return
	}

	passwordOld, err := db.Client.GetById(int(passwordId), int(userId))
	if err != nil {
		error_tracer.Client.InfoLog("updatePassword", "notFound",
			fmt.Sprintf("userId %d passwordId %d not found Error : %v", userId, passwordId, err))
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Password not found",
			"data":    err.Error(),
		})

		return
	}

	passwordOld.AccountName = password.AccountName
	passwordOld.Username = password.Username
	passwordOld.Email = password.Email
	passwordOld.Password = password.Password

	if err := db.Client.Save(&passwordOld); err != nil {
		error_tracer.Client.ErrorLog("updatePassword", "database", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "can't update in database",
			"data":    err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successful",
	})
}
