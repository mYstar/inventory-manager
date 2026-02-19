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
	var responseItem storage.Item
	err := json.Unmarshal(response.Body.Bytes(), &responseItem)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, expected, responseItem)
}

func TestSubtractFromQuantity(t *testing.T) {
	router := initTestRouter()
	request := DeltaRequest{QuantityDelta: new(int64(-3))}
	requestJson, _ := json.Marshal(request)

	response := sendTestRequest(router, "PATCH", "/item/2", string(requestJson))

	expected := storage.Item{Name: "Chair", PriceCents: 2489, Quantity: 1}
	var responseItem storage.Item
	err := json.Unmarshal(response.Body.Bytes(), &responseItem)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, expected, responseItem)
}

func TestSubtractUnderflow(t *testing.T) {
	router := initTestRouter()
	request := DeltaRequest{QuantityDelta: new(int64(-7))}
	requestJson, _ := json.Marshal(request)

	response := sendTestRequest(router, "PATCH", "/item/1", string(requestJson))

	expected := NewError("Current quantity is too small to perform the operation.")
	var responseItem ErrorResponse
	err := json.Unmarshal(response.Body.Bytes(), &responseItem)
	assert.Nil(t, err)
	assert.Equal(t, 409, response.Code)
	assert.Equal(t, expected, responseItem)
}

func TestAlterUnknownItem(t *testing.T) {
	router := initTestRouterEmpty()
	request := DeltaRequest{QuantityDelta: new(int64(-7))}
	requestJson, _ := json.Marshal(request)
	response := sendTestRequest(router, "PATCH", "/item/1", string(requestJson))

	expected := NewError("Item ID does not exist.")
	var responseItem ErrorResponse
	err := json.Unmarshal(response.Body.Bytes(), &responseItem)
	assert.Nil(t, err)
	assert.Equal(t, 404, response.Code)
	assert.Equal(t, expected, responseItem)
}

func TestAlterInvalidId(t *testing.T) {
	router := initTestRouterEmpty()
	response := sendTestRequest(router, "PATCH", "/item/unknown", "{}")

	expected := NewError("Invalid item ID.")
	var responseItem ErrorResponse
	err := json.Unmarshal(response.Body.Bytes(), &responseItem)
	assert.Nil(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, expected, responseItem)
}

func TestAlterInvalidBody(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "PATCH", "/item/1", "{")

	expected := NewError("Request body is not valid JSON.")
	var responseItem ErrorResponse
	err := json.Unmarshal(response.Body.Bytes(), &responseItem)
	assert.Nil(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, expected, responseItem)
}

func TestAlterMissingDelta(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "PATCH", "/item/1", "{}")

	expected := NewError("Request body does not contain required fields.")
	var responseItem ErrorResponse
	err := json.Unmarshal(response.Body.Bytes(), &responseItem)
	assert.Nil(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, expected, responseItem)
}
