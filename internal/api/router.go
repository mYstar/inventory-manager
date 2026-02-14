package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/items", getItems)
	router.POST("/item", createItem)
	router.DELETE("/item/:id", deleteItem)
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
		c.JSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}
	isInvalid := item.Name == "" || item.Price == 0 || item.Quantity == 0
	if isInvalid {
		c.JSON(http.StatusBadRequest, NewError("Missing item data"))
		return
	}

	Items[uint(len(Items))+1] = item

	c.JSON(http.StatusCreated, item)
}

// deleteItem deletes an item from the inventory.
func deleteItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewError("Invalid item ID."))
		return
	}
	_, exists := Items[uint(id)]
	if !exists {
		c.JSON(http.StatusNotFound, NewError("Item ID does not exist."))
		return
	}

	delete(Items, uint(id))

	c.JSON(http.StatusNoContent, nil)
}

func NewError(message string) map[string]string {
	return map[string]string{"error": message}
}
