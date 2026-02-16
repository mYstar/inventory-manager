package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeltaRequest struct {
	QuantityDelta *int `json:"quantity_delta"`
}

type IdsQuery struct {
	Ids []uint `form:"ids"`
}

type ItemResponse struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity uint    `json:"quantity"`
}

type ValueResponse struct {
	Success bool    `json:"success"`
	Value   float32 `json:"value"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/items", getItems)
	router.GET("/items/value", getItemsValue)
	router.POST("/item", createItem)
	router.PATCH("/item/:id", alterQuantity)
	router.DELETE("/item/:id", deleteItem)
	return router
}

// getItems responds with the list of all albums as JSON.
func getItems(c *gin.Context) {
	c.JSON(http.StatusOK, Inventory)
}

// getItemsValue responds with the value sum of the given Items or all Items if no Item IDs are given in the request body.
func getItemsValue(c *gin.Context) {
	query := IdsQuery{}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, NewError("One or more of given ids cannot be interpreted as uint."))
		return
	}

	value := Inventory.calculateValue(query.Ids)

	response := ValueResponse{Success: true, Value: value}
	c.JSON(http.StatusOK, response)
}

// createItem creates a new item and adds it to the inventory.
func createItem(c *gin.Context) {
	var item Item
	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewError("Request body is not valid JSON."))
		return
	}
	isInvalid := item.Name == "" || item.Price == 0 || item.Quantity == 0
	if isInvalid {
		c.JSON(http.StatusBadRequest, NewError("Missing item data."))
		return
	}

	newIdx, newItem := Inventory.createItem(item)

	response := ItemResponse{ID: newIdx, Name: newItem.Name, Price: newItem.Price, Quantity: newItem.Quantity}
	c.JSON(http.StatusCreated, response)
}

// deleteItem deletes an item from the inventory.
func deleteItem(c *gin.Context) {
	id, httpCode, errResponse := getId(c)
	if errResponse != nil {
		c.JSON(httpCode, errResponse)
		return
	}

	Inventory.delete(id)

	c.JSON(http.StatusNoContent, nil)
}

// alterQuantity deletes an item from the inventory.
func alterQuantity(c *gin.Context) {
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
	item, errResponse := Inventory.alterQuantity(id, dQuantity)
	if errResponse != nil {
		c.JSON(http.StatusConflict, errResponse)
		return
	}

	c.JSON(http.StatusOK, item)
}

// getId extracts and validates the "id" parameter from the request.
// returns the item ID or an error if the id can't be parsed or does not exist.
func getId(c *gin.Context) (uint, int, *ErrorResponse) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, http.StatusBadRequest, new(NewError("Invalid item ID."))
	}
	exists := Inventory.itemExists(uint(id))
	if !exists {
		return 0, http.StatusNotFound, new(NewError("Item ID does not exist."))
	}
	return uint(id), 0, nil
}

func NewError(message string) ErrorResponse {
	return ErrorResponse{Error: message}
}
