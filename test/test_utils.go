package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"inventory_manager/internal/db"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// initTestRouter initializes the test router and fills the inventory with sample data.
func initTestRouter() *gin.Engine {
	inventory := api.NewInventory(db.NewMemoryPersistence())
	router := gin.Default()
	api.SetupRoutes(router, inventory)

	sendTestRequest(router, "POST", "/item", `{"name": "Book", "price_cents": 1099, "quantity": 1}`)
	sendTestRequest(router, "POST", "/item", `{"name": "Chair", "price_cents": 2489, "quantity": 4}`)
	sendTestRequest(router, "POST", "/item", `{"name": "Laptop", "price_cents": 109900, "quantity": 3}`)

	return router
}

// initTestRouterEmpty initializes a router with an empty inventory.
func initTestRouterEmpty() *gin.Engine {
	inventory := api.NewInventory(db.NewMemoryPersistence())
	router := gin.Default()
	api.SetupRoutes(router, inventory)

	return router
}

// sendTestRequest sends a request to the router and returns the response.
func sendTestRequest(router *gin.Engine, method, url, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	router.ServeHTTP(w, req)

	return w
}

// assertInventoryEquals checks if the inventory of the router matches the expected items.
func assertInventoryEquals(t *testing.T, router *gin.Engine, expectedItems db.Items) {
	response := sendTestRequest(router, "GET", "/items", "")
	expectedJson, _ := json.Marshal(expectedItems)

	assert.Equal(t, 200, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}
