package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateItemsValue(t *testing.T) {
	originalItems := maps.Clone(api.Inventory)
	t.Cleanup(func() { api.Inventory = originalItems })

	testItems := make(map[uint]api.Item)
	testItems[1] = api.Item{Name: "Book", Price: 10.99, Quantity: 1}
	testItems[2] = api.Item{Name: "Chair", Price: 24.89, Quantity: 4}
	testItems[3] = api.Item{Name: "Laptop", Price: 1099.00, Quantity: 3}
	api.Inventory = testItems

	response := sendTestRequest("GET", "/items/value?ids=1&ids=3", "")

	value := api.ValueResponse{}
	err := json.Unmarshal([]byte(response.Body.String()), &value)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, true, value.Success)
	assert.Equal(t, float32(3307.99), value.Value)
}

func TestCalculateMissingIds(t *testing.T) {
	response := sendTestRequest("GET", "/items/value?ids=1&ids=2&ids=3", "")

	value := api.ValueResponse{}
	err := json.Unmarshal([]byte(response.Body.String()), &value)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, true, value.Success)
	assert.Equal(t, float32(0.0), value.Value)
}

func TestCalculateInvalidIds(t *testing.T) {
	response := sendTestRequest("GET", "/items/value?ids=a&ids=2&ids=c", "")

	value := api.ErrorResponse{}
	err := json.Unmarshal([]byte(response.Body.String()), &value)

	assert.Nil(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, false, value.Success)
	assert.Equal(t, "One or more of given ids cannot be interpreted as uint.", value.Error)
}
