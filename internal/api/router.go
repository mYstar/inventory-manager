package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Delta struct {
	QuantityDelta int
}

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/items", getItems)
	router.POST("/item", createItem)
	router.PATCH("/item/:id", alterQuantity)
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
	id, err := getId(c)
	if err != nil {
		return
	}

	delete(Items, id)

	c.JSON(http.StatusNoContent, nil)
}

// alterQuantity deletes an item from the inventory.
func alterQuantity(c *gin.Context) {
	id, err := getId(c)
	if err != nil {
		return
	}

	var operation Delta
	err = c.BindJSON(&operation)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewError(err.Error()))
		return
	}

	var dQuantity = operation.QuantityDelta
	var item = Items[id]
	switch {
	case dQuantity < 0:
		if uint(-dQuantity) > item.Quantity {
			c.JSON(http.StatusConflict, NewError("Quantity is too small to perform the operation."))
			return
		}
		item.Quantity -= uint(-dQuantity)
	default:
		item.Quantity += uint(dQuantity)
	}

	c.JSON(http.StatusOK, item)
}

// getId extracts and validates the "id" parameter from the request.
// returns the item ID or an error if the id can't be parsed or does not exist.
func getId(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewError("Invalid item ID."))
		return 0, err
	}
	_, exists := Items[uint(id)]
	if !exists {
		c.JSON(http.StatusNotFound, NewError("Item ID does not exist."))
		return 0, errors.New("item id does not exist")
	}
	return uint(id), nil
}

func NewError(message string) map[string]string {
	return map[string]string{"error": message}
}
