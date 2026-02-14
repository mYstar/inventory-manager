package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetItems(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })

	testItems := make(map[uint]api.Item)
	testItems[1] = api.Item{Name: "Book", Price: 10.99, Quantity: 1}
	testItems[2] = api.Item{Name: "Chair", Price: 24.89, Quantity: 4}
	testItems[3] = api.Item{Name: "Laptop", Price: 1099.00, Quantity: 3}
	api.Items = testItems

	response := sendTestRequest("GET", "/items", "")

	expected, _ := json.Marshal(api.Items)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, string(expected), response.Body.String())
}

func TestGetEmptyItems(t *testing.T) {
	response := sendTestRequest("GET", "/items", "")

	expected := "{}"
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, expected, response.Body.String())
}
