package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddToQuantity(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })
	api.Items[1] = api.Item{Name: "Book", Price: 12.99, Quantity: 1}
	requestDelta := api.Delta{QuantityDelta: 3}
	requestDeltaJson, _ := json.Marshal(requestDelta)

	response := sendTestRequest("PATCH", "/item/1", string(requestDeltaJson))

	expected := api.Item{Name: "Book", Price: 12.99, Quantity: 4}
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestSubtractFromQuantity(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })
	api.Items[1] = api.Item{Name: "Book", Price: 12.99, Quantity: 5}
	requestDelta := api.Delta{QuantityDelta: -3}
	requestDeltaJson, _ := json.Marshal(requestDelta)

	response := sendTestRequest("PATCH", "/item/1", string(requestDeltaJson))

	expected := api.Item{Name: "Book", Price: 12.99, Quantity: 2}
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestSubtractUnderflow(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })
	api.Items[1] = api.Item{Name: "Book", Price: 12.99, Quantity: 5}
	requestDelta := api.Delta{QuantityDelta: -7}
	requestDeltaJson, _ := json.Marshal(requestDelta)

	response := sendTestRequest("PATCH", "/item/1", string(requestDeltaJson))

	expected := api.NewError("Quantity is too small to perform the operation.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 409, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestAlterUnknownItem(t *testing.T) {
	response := sendTestRequest("PATCH", "/item/1", "{}")

	expected := api.NewError("Item ID does not exist.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 404, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestAlterInvalidId(t *testing.T) {
	response := sendTestRequest("PATCH", "/item/unknown", "{}")

	expected := api.NewError("Invalid item ID.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}
