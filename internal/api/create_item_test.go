package api

import (
	"encoding/json"
	"inventory_manager/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateItem(t *testing.T) {
	router := initTestRouterEmpty()
	testItem := storage.Item{
		Name:       "Book",
		PriceCents: 1099,
		Quantity:   1,
	}
	testItemJson, _ := json.Marshal(testItem)

	response := sendTestRequest(router, "POST", "/item", string(testItemJson))

	createdItemJson := response.Body.String()
	createdItem := ItemResponse{}
	err := json.Unmarshal([]byte(createdItemJson), &createdItem)

	assert.Nil(t, err)
	assert.Equal(t, 201, response.Code)
	assert.Equal(t, uint64(1), createdItem.ID)
	assert.Equal(t, testItem.Name, createdItem.Name)
	assert.Equal(t, testItem.Quantity, createdItem.Quantity)
	assert.Equal(t, testItem.PriceCents, createdItem.PriceCents)

	// check if the new item is in the inventory
	expected := storage.Items{1: testItem}
	assertInventoryEquals(t, router, expected)
}

func TestCreateItemExistingValues(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "DELETE", "/item/2", "")
	testItem := storage.Item{
		Name:       "Sofa",
		PriceCents: 50999,
		Quantity:   3,
	}
	testItemJson, _ := json.Marshal(testItem)

	// check if creating works
	response = sendTestRequest(router, "POST", "/item", string(testItemJson))

	createdItemJson := response.Body.String()
	createdItem := ItemResponse{}
	err := json.Unmarshal([]byte(createdItemJson), &createdItem)

	assert.Nil(t, err)
	assert.Equal(t, 201, response.Code)
	assert.Equal(t, uint64(4), createdItem.ID)
	assert.Equal(t, testItem.Name, createdItem.Name)
	assert.Equal(t, testItem.Quantity, createdItem.Quantity)
	assert.Equal(t, testItem.PriceCents, createdItem.PriceCents)

	// check if the new item is in the inventory
	expected := storage.Items{
		1: storage.Item{Name: "Book", PriceCents: 1099, Quantity: 1},
		3: storage.Item{Name: "Laptop", PriceCents: 109900, Quantity: 3},
		4: testItem,
	}
	assertInventoryEquals(t, router, expected)
}

func TestCreateItemWithInvalidJson(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "POST", "/item", "{")

	expected := NewError("Request body invalid JSON or is missing required fields.")
	var responseItem ErrorResponse
	err := json.Unmarshal(response.Body.Bytes(), &responseItem)
	assert.Nil(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, expected, responseItem)
}

func TestCreateItemWithMissingValues(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "POST", "/item", "{\"name\": \"Book\", \"price_cents\": 1099}")

	expected := NewError("Request body invalid JSON or is missing required fields.")
	var responseItem ErrorResponse
	err := json.Unmarshal(response.Body.Bytes(), &responseItem)
	assert.Nil(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, expected, responseItem)
}

func TestCreateItemWithNegativeQuantity(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "POST", "/item", "{\"name\": \"Book\", \"price_cents\": 1099, \"quantity\": -1}")

	expected := NewError("Quantity cannot be negative.")
	var responseItem ErrorResponse
	err := json.Unmarshal(response.Body.Bytes(), &responseItem)
	assert.Nil(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, expected, responseItem)
}
