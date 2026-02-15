package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTotalValue(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })

	testItems := make(map[uint]api.Item)
	testItems[1] = api.Item{Name: "Book", Price: 10.99, Quantity: 1}
	testItems[2] = api.Item{Name: "Chair", Price: 24.89, Quantity: 4}
	testItems[3] = api.Item{Name: "Laptop", Price: 1099.00, Quantity: 3}
	api.Items = testItems

	response := sendTestRequest("GET", "/items/value", "")

	value := api.ValueResponse{}
	err := json.Unmarshal([]byte(response.Body.String()), &value)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, true, value.Success)
	assert.Equal(t, float32(3407.55), value.Value)
}

func TestCalculateTotalValueEmpty(t *testing.T) {
	response := sendTestRequest("GET", "/items/value", "")

	value := api.ValueResponse{}
	err := json.Unmarshal([]byte(response.Body.String()), &value)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, true, value.Success)
	assert.Equal(t, float32(0.0), value.Value)
}
