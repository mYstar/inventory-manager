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
	request := api.DeltaRequest{QuantityDelta: new(3)}
	requestJson, _ := json.Marshal(request)

	response := sendTestRequest("PATCH", "/item/1", string(requestJson))

	expected := api.Item{Name: "Book", Price: 12.99, Quantity: 4}
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestSubtractFromQuantity(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })
	api.Items[1] = api.Item{Name: "Book", Price: 12.99, Quantity: 5}
	request := api.DeltaRequest{QuantityDelta: new(-3)}
	requestJson, _ := json.Marshal(request)

	response := sendTestRequest("PATCH", "/item/1", string(requestJson))

	expected := api.Item{Name: "Book", Price: 12.99, Quantity: 2}
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestSubtractUnderflow(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })
	api.Items[1] = api.Item{Name: "Book", Price: 12.99, Quantity: 5}
	request := api.DeltaRequest{QuantityDelta: new(-7)}
	requestJson, _ := json.Marshal(request)

	response := sendTestRequest("PATCH", "/item/1", string(requestJson))

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

func TestAlterInvalidBody(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })
	api.Items[1] = api.Item{Name: "Book", Price: 12.99, Quantity: 1}

	response := sendTestRequest("PATCH", "/item/1", "{")

	expected := api.NewError("Request body is not valid JSON.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestAlterMissingDelta(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })
	api.Items[1] = api.Item{Name: "Book", Price: 12.99, Quantity: 1}

	response := sendTestRequest("PATCH", "/item/1", "{}")

	expected := api.NewError("Request body does not contain required fields.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}
