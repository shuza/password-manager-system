package api

import (
	"github.com/gin-gonic/gin"
)

func NewGinEngine() *gin.Engine {
	r := gin.Default()
	r.GET("/", Index)

	routes := r.Group("/api/v1")
	routes.POST("/password", addPassword)
	routes.GET("/password", passwordList)

	routes.GET("/password/:password_id", passwordDetails)
	routes.PUT("/password/:password_id", updatePassword)
	routes.DELETE("/password/:password_id", deletePassword)

	return r
}
