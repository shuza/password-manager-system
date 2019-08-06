package api

import (
	"github.com/gin-gonic/gin"
)

func NewGinEngine() *gin.Engine {
	r := gin.Default()
	r.GET("/", Index)

	//routes := r.Group("/api/v1")

	return r
}
