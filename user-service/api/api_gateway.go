package api

import (
	"github.com/gin-gonic/gin"
)

func NewGinEngine() *gin.Engine {
	r := gin.Default()
	r.GET("/user", index)

	routes := r.Group("/user/api/v1")
	routes.POST("/user", createUser)
	routes.POST("/auth/login", login)
	routes.GET("/auth/token", tokenVerify)

	return r
}
