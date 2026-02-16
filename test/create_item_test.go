package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateItem(t *testing.T) {
	originalItems := maps.Clone(api.Inventory)
	t.Cleanup(func() { api.Inventory = originalItems })
	testItem := api.Item{
		Name:     "Book",
		Price:    10.99,
		Quantity: 1,
	}
	testItemJson, _ := json.Marshal(testItem)

	response := sendTestRequest("POST", "/item", string(testItemJson))

	createdItemJson := response.Body.String()
	createdItem := api.ItemResponse{}
	err := json.Unmarshal([]byte(createdItemJson), &createdItem)

	assert.Nil(t, err)
	assert.Equal(t, 201, response.Code)
	assert.Equal(t, uint(1), createdItem.ID)
	assert.Equal(t, testItem.Name, createdItem.Name)
	assert.Equal(t, testItem.Quantity, createdItem.Quantity)
	assert.Equal(t, testItem.Price, createdItem.Price)
}

func TestCreateItemExistingValues(t *testing.T) {
	originalItems := maps.Clone(api.Inventory)
	t.Cleanup(func() { api.Inventory = originalItems })

	initialItems := make(map[uint]api.Item)
	initialItems[1] = api.Item{Name: "Book", Price: 10.99, Quantity: 1}
	initialItems[3] = api.Item{Name: "Laptop", Price: 1099.00, Quantity: 3}
	api.Inventory = maps.Clone(initialItems)

	testItem := api.Item{
		Name:     "Sofa",
		Price:    509.99,
		Quantity: 3,
	}
	testItemJson, _ := json.Marshal(testItem)

	// check if creating works
	response := sendTestRequest("POST", "/item", string(testItemJson))

	createdItemJson := response.Body.String()
	createdItem := api.ItemResponse{}
	err := json.Unmarshal([]byte(createdItemJson), &createdItem)

	assert.Nil(t, err)
	assert.Equal(t, 201, response.Code)
	assert.Equal(t, uint(4), createdItem.ID)
	assert.Equal(t, testItem.Name, createdItem.Name)
	assert.Equal(t, testItem.Quantity, createdItem.Quantity)
	assert.Equal(t, testItem.Price, createdItem.Price)

	// check if the new item is in the inventory
	response = sendTestRequest("GET", "/items", "")

	initialItems[4] = testItem
	expected, _ := json.Marshal(initialItems)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, string(expected), response.Body.String())
}

func TestCreateItemWithInvalidJson(t *testing.T) {
	response := sendTestRequest("POST", "/item", "{")

	expected := api.NewError("Request body is not valid JSON.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestCreateItemWithMissingValues(t *testing.T) {
	response := sendTestRequest("POST", "/item", "{\"name\": \"Book\", \"price\": 10.99}")

	expected := api.NewError("Missing item data.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}
