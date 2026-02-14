package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/items", getItems)
	router.POST("/item", createItem)
	return router
}

// getItems responds with the list of all albums as JSON.
func getItems(c *gin.Context) {
	c.JSON(http.StatusOK, Items)
}

// createItem creates a new item and adds it to the inventory.
func createItem(c *gin.Context) {
	var item Item

	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isInvalid := item.Name == "" || item.Price == 0 || item.Quantity == 0
	if isInvalid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing item data"})
		return
	}

	item.ID = uint(len(Items)) + 1 //TODO: calculate real max value
	Items = append(Items, item)

	c.JSON(http.StatusCreated, item)
}
