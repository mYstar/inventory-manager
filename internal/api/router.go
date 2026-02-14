package api

import "github.com/gin-gonic/gin"

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/items", getItems)
	return router
}
