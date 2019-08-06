package api

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/db"
	"user-service/error_tracer"
	"user-service/model"
	"user-service/service"
)

func createUser(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		error_tracer.Client.ErrorLog("createUser", "requestBody", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"data":    err.Error(),
		})

		return
	}

	if user.Email == "" || user.Password == "" || user.Fullname == "" || user.BusinessName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"data":    "Required fields can't be empty",
		})

		return
	}

	h := md5.New()
	h.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(h.Sum(nil))

	if err := db.Client.Save(&user); err != nil {
		error_tracer.Client.ErrorLog("createUser", "database", err.Error())
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

func login(c *gin.Context) {
	var reqUser model.User
	if err := c.BindJSON(&reqUser); err != nil {
		error_tracer.Client.ErrorLog("login", "requestBody", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"data":    err.Error(),
		})

		return
	}

	user, err := db.Client.GetByEmail(reqUser.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
			"data":    err.Error(),
		})

		return
	}

	h := md5.New()
	h.Write([]byte(reqUser.Password))
	reqUser.Password = hex.EncodeToString(h.Sum(nil))

	if user.Password != reqUser.Password {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Invalid credential",
		})

		return
	}

	tokenService := service.TokenService{}
	token, err := tokenService.Encode(user)
	if err != nil {
		error_tracer.Client.ErrorLog("login", "token", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Can't generate token",
			"data":    err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    token,
	})
}
