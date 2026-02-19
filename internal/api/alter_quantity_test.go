package api

import (
	"encoding/json"
	"inventory_manager/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddToQuantity(t *testing.T) {
	router := initTestRouter()
	request := DeltaRequest{QuantityDelta: new(int64(3))}
	requestJson, _ := json.Marshal(request)

	response := sendTestRequest(router, "PATCH", "/item/1", string(requestJson))

	expected := storage.Item{Name: "Book", PriceCents: 1099, Quantity: 4}
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestSubtractFromQuantity(t *testing.T) {
	router := initTestRouter()
	request := DeltaRequest{QuantityDelta: new(int64(-3))}
	requestJson, _ := json.Marshal(request)

	response := sendTestRequest(router, "PATCH", "/item/2", string(requestJson))

	expected := storage.Item{Name: "Chair", PriceCents: 2489, Quantity: 1}
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestSubtractUnderflow(t *testing.T) {
	router := initTestRouter()
	request := DeltaRequest{QuantityDelta: new(int64(-7))}
	requestJson, _ := json.Marshal(request)

	response := sendTestRequest(router, "PATCH", "/item/1", string(requestJson))

	expected := NewError("Current quantity is too small to perform the operation.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 409, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestAlterUnknownItem(t *testing.T) {
	router := initTestRouterEmpty()
	request := DeltaRequest{QuantityDelta: new(int64(-7))}
	requestJson, _ := json.Marshal(request)
	response := sendTestRequest(router, "PATCH", "/item/1", string(requestJson))

	expected := NewError("Item ID does not exist.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 404, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestAlterInvalidId(t *testing.T) {
	router := initTestRouterEmpty()
	response := sendTestRequest(router, "PATCH", "/item/unknown", "{}")

	expected := NewError("Invalid item ID.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestAlterInvalidBody(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "PATCH", "/item/1", "{")

	expected := NewError("Request body is not valid JSON.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestAlterMissingDelta(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "PATCH", "/item/1", "{}")

	expected := NewError("Request body does not contain required fields.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}
