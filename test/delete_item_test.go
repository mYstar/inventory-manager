package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"inventory_manager/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteItem(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "DELETE", "/item/1", "")

	assert.Equal(t, 204, response.Code)
	assert.Equal(t, "", response.Body.String())

	// test if item is actually deleted
	expected := db.Items{
		2: db.Item{Name: "Chair", Price: 24.89, Quantity: 4},
		3: db.Item{Name: "Laptop", Price: 1099.00, Quantity: 3},
	}
	assertInventoryEquals(t, router, expected)
}

func TestDeleteUnknownItem(t *testing.T) {
	router := initTestRouterEmpty()
	response := sendTestRequest(router, "DELETE", "/item/1", "")

	expected := api.NewError("Item ID does not exist.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 404, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestDeleteInvalidId(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "DELETE", "/item/invalid", "")

	expected := api.NewError("Invalid item ID.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}
