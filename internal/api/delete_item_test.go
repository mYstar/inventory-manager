package api

import (
	"encoding/json"
	"inventory_manager/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteItem(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "DELETE", "/item/1", "")

	assert.Equal(t, 204, response.Code)
	assert.Equal(t, "", response.Body.String())

	// test if item is actually deleted
	expected := storage.Items{
		2: storage.Item{Name: "Chair", PriceCents: 2489, Quantity: 4},
		3: storage.Item{Name: "Laptop", PriceCents: 109900, Quantity: 3},
	}
	assertInventoryEquals(t, router, expected)
}

func TestDeleteUnknownItem(t *testing.T) {
	router := initTestRouterEmpty()
	response := sendTestRequest(router, "DELETE", "/item/1", "")

	expected := NewError("Item ID does not exist.")
	var responseItem ErrorResponse
	err := json.Unmarshal(response.Body.Bytes(), &responseItem)
	assert.Nil(t, err)
	assert.Equal(t, 404, response.Code)
	assert.Equal(t, expected, responseItem)
}

func TestDeleteInvalidId(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "DELETE", "/item/invalid", "")

	expected := NewError("Invalid item ID.")
	var responseItem ErrorResponse
	err := json.Unmarshal(response.Body.Bytes(), &responseItem)
	assert.Nil(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, expected, responseItem)
}
