package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"inventory_manager/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddToQuantity(t *testing.T) {
	router := initTestRouter()
	request := api.DeltaRequest{QuantityDelta: new(int64(3))}
	requestJson, _ := json.Marshal(request)

	response := sendTestRequest(router, "PATCH", "/item/1", string(requestJson))

	expected := db.Item{Name: "Book", PriceCents: 1099, Quantity: 4}
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestSubtractFromQuantity(t *testing.T) {
	router := initTestRouter()
	request := api.DeltaRequest{QuantityDelta: new(int64(-3))}
	requestJson, _ := json.Marshal(request)

	response := sendTestRequest(router, "PATCH", "/item/2", string(requestJson))

	expected := db.Item{Name: "Chair", PriceCents: 2489, Quantity: 1}
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestSubtractUnderflow(t *testing.T) {
	router := initTestRouter()
	request := api.DeltaRequest{QuantityDelta: new(int64(-7))}
	requestJson, _ := json.Marshal(request)

	response := sendTestRequest(router, "PATCH", "/item/1", string(requestJson))

	expected := api.NewError("Quantity is too small to perform the operation.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 409, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestAlterUnknownItem(t *testing.T) {
	router := initTestRouterEmpty()
	response := sendTestRequest(router, "PATCH", "/item/1", "{}")

	expected := api.NewError("Item ID does not exist.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 404, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestAlterInvalidId(t *testing.T) {
	router := initTestRouterEmpty()
	response := sendTestRequest(router, "PATCH", "/item/unknown", "{}")

	expected := api.NewError("Invalid item ID.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestAlterInvalidBody(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "PATCH", "/item/1", "{")

	expected := api.NewError("Request body is not valid JSON.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestAlterMissingDelta(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "PATCH", "/item/1", "{}")

	expected := api.NewError("Request body does not contain required fields.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}
