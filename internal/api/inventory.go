package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Item struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity uint    `json:"quantity"`
}

var Items = make([]Item, 0)

// getItems responds with the list of all albums as JSON.
func getItems(c *gin.Context) {
	c.JSON(http.StatusOK, Items)
}
