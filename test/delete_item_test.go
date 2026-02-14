package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteItem(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })
	api.Items[1] = api.Item{Name: "Book", Price: 12.99, Quantity: 1}

	response := sendTestRequest("DELETE", "/item/1", "")

	assert.Equal(t, 204, response.Code)
	assert.Equal(t, "", response.Body.String())
}

func TestDeleteUnknownItem(t *testing.T) {
	response := sendTestRequest("DELETE", "/item/1", "")

	expected := api.NewError("Item ID does not exist.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 404, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestDeleteInvalidId(t *testing.T) {
	response := sendTestRequest("DELETE", "/item/invalid", "")

	expected := api.NewError("Invalid item ID.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}
