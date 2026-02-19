package api

import (
	"inventory_manager/internal/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetupRoutes defines API endpoints for inventory management, attaching handlers for inventory operations.
func SetupRoutes(router *gin.Engine, inventory *Inventory) {
	router.GET("/items", func(c *gin.Context) { getItems(c, inventory) })
	router.GET("/items/value", func(c *gin.Context) { getItemsValue(c, inventory) })
	router.POST("/item", func(c *gin.Context) { createItem(c, inventory) })
	router.PATCH("/item/:id", func(c *gin.Context) { alterQuantity(c, inventory) })
	router.DELETE("/item/:id", func(c *gin.Context) { deleteItem(c, inventory) })
}

// getItems responds with the list of all items in the inventory as JSON (type storage.Items).
func getItems(c *gin.Context, inventory *Inventory) {
	items, err := inventory.GetItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError("Internal server error: "+err.Error()))
		return
	}
	c.JSON(http.StatusOK, items)
}

// getItemsValue responds with the total value of a list of Items (type ValueResponse).
// If a list of ids is given via an URL query the items of this list are summed up,
// otherwise all Items in the inventory are used.
func getItemsValue(c *gin.Context, inventory *Inventory) {
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

// createItem creates a new item from the values in the request body (type storage.Item) and adds it to the inventory.
func createItem(c *gin.Context, inventory *Inventory) {
	var item storage.Item
	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewError("Request body invalid JSON or is missing required fields."))
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

// deleteItem deletes the item with the ID provided in the URL from the inventory.
func deleteItem(c *gin.Context, inventory *Inventory) {
	id, httpCode, errResponse := getId(c)
	if errResponse != nil {
		c.JSON(httpCode, errResponse)
		return
	}

	err := inventory.delete(id)
	if err != nil {
		switch err.Error() {
		case "item id does not exist":
			c.JSON(http.StatusNotFound, NewError("Item ID does not exist."))
		default:
			c.JSON(http.StatusInternalServerError, NewError("Internal server error: "+err.Error()))
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// alterQuantity adjusts the quantity of an inventory item with the ID provided in the URL
// and delta from the HTTP request body (type DeltaRequest).
func alterQuantity(c *gin.Context, inventory *Inventory) {
	id, httpCode, errResponse := getId(c)
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
		switch err.Error() {
		case "item id does not exist":
			c.JSON(http.StatusNotFound, NewError("Item ID does not exist."))
		case "quantity is too small to perform the operation":
			c.JSON(http.StatusConflict, NewError("Current quantity is too small to perform the operation."))
		default:
			c.JSON(http.StatusInternalServerError, NewError("Internal server error: "+err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, item)
}

// getId extracts and validates the ID parameter from the URL.
// returns the item ID or an error if the id can't be parsed or does not exist.
func getId(c *gin.Context) (uint, int, *ErrorResponse) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, http.StatusBadRequest, &ErrorResponse{Success: false, Error: "Invalid item ID."}
	}
	return uint(id), 0, nil
}

// NewError creates a new ErrorResponse with the given message.
func NewError(message string) ErrorResponse {
	return ErrorResponse{Success: false, Error: message}
}
