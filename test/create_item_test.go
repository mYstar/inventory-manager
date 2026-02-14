package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateItem(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })
	testItem := api.Item{
		Name:     "Book",
		Price:    10.99,
		Quantity: 1,
	}
	testItemJson, _ := json.Marshal(testItem)

	response := sendTestRequest("POST", "/item", string(testItemJson))

	createdItemJson := response.Body.String()
	createdItem := api.Item{}
	err := json.Unmarshal([]byte(createdItemJson), &createdItem)

	assert.Nil(t, err)
	assert.Equal(t, 201, response.Code)
	assert.Equal(t, testItem.Name, createdItem.Name)
	assert.Equal(t, testItem.Quantity, createdItem.Quantity)
	assert.Equal(t, testItem.Price, createdItem.Price)
}

func TestCreateItemWithInvalidJson(t *testing.T) {
	response := sendTestRequest("POST", "/item", "{")

	assert.Equal(t, 400, response.Code)
}

func TestCreateItemWithMissingValues(t *testing.T) {
	response := sendTestRequest("POST", "/item", "{\"name\": \"Book\", \"price\": 10.99}")

	assert.Equal(t, 400, response.Code)
}
