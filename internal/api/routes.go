package api

import (
	"inventory_manager/internal/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, inventory Inventory) {
	router.GET("/items", func(c *gin.Context) { getItems(c, inventory) })
	router.GET("/items/value", func(c *gin.Context) { getItemsValue(c, inventory) })
	router.POST("/item", func(c *gin.Context) { createItem(c, inventory) })
	router.PATCH("/item/:id", func(c *gin.Context) { alterQuantity(c, inventory) })
	router.DELETE("/item/:id", func(c *gin.Context) { deleteItem(c, inventory) })
}

// getItems responds with the list of all albums as JSON.
func getItems(c *gin.Context, inventory Inventory) {
	items, err := inventory.persistence.GetItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError("Internal server error: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, items)
}

// getItemsValue responds with the value sum of the given Items or all Items if no Item IDs are given in the request body.
func getItemsValue(c *gin.Context, inventory Inventory) {
	query := IdsQuery{}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, NewError("One or more of given ids cannot be interpreted as uint."))
		return
	}

	value, err := inventory.calculateValue(query.Ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError("Internal server error: "+err.Error()))
		return
	}

	response := ValueResponse{Success: true, ValueCents: value}
	c.JSON(http.StatusOK, response)
}

// createItem creates a new item and adds it to the inventory.
func createItem(c *gin.Context, inventory Inventory) {
	var item db.Item
	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewError("Request body is not valid JSON."))
		return
	}
	isInvalid := item.Name == "" || item.PriceCents == 0 || item.Quantity == 0
	if isInvalid {
		c.JSON(http.StatusBadRequest, NewError("Missing item data."))
		return
	}
	if item.Quantity < 0 {
		c.JSON(http.StatusBadRequest, NewError("Quantity cannot be negative."))
		return
	}

	newIdx, newItem, err := inventory.createItem(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError("Internal server error: "+err.Error()))
		return
	}

	response := ItemResponse{ID: newIdx, Name: newItem.Name, PriceCents: newItem.PriceCents, Quantity: newItem.Quantity}
	c.JSON(http.StatusCreated, response)
}

// deleteItem deletes an item from the inventory.
func deleteItem(c *gin.Context, inventory Inventory) {
	id, httpCode, errResponse := getId(c, inventory)
	if errResponse != nil {
		c.JSON(httpCode, errResponse)
		return
	}

	err := inventory.delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError("Internal server error: "+err.Error()))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// alterQuantity deletes an item from the inventory.
func alterQuantity(c *gin.Context, inventory Inventory) {
	id, httpCode, errResponse := getId(c, inventory)
	if errResponse != nil {
		c.JSON(httpCode, errResponse)
		return
	}

	var request DeltaRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, NewError("Request body is not valid JSON."))
		return
	}
	if request.QuantityDelta == nil {
		c.JSON(http.StatusBadRequest, NewError("Request body does not contain required fields."))
		return
	}

	var dQuantity = *request.QuantityDelta
	item, err := inventory.alterQuantity(id, dQuantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError("Internal server error: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, item)
}

// getId extracts and validates the "id" parameter from the request.
// returns the item ID or an error if the id can't be parsed or does not exist.
func getId(c *gin.Context, inventory Inventory) (uint, int, *ErrorResponse) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, http.StatusBadRequest, new(NewError("Invalid item ID."))
	}
	exists := inventory.itemExists(uint(id))
	if !exists {
		return 0, http.StatusNotFound, new(NewError("Item ID does not exist."))
	}
	return uint(id), 0, nil
}

func NewError(message string) ErrorResponse {
	return ErrorResponse{Success: false, Error: message}
}
