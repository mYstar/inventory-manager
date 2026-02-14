package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddToQuantity(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })
	api.Items[1] = api.Item{Name: "Book", Price: 12.99, Quantity: 1}
	requestDelta := api.Delta{QuantityDelta: 3}
	requestDeltaJson, _ := json.Marshal(requestDelta)

	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/item/1", strings.NewReader(string(requestDeltaJson)))
	router.ServeHTTP(w, req)

	expected := api.Item{Name: "Book", Price: 12.99, Quantity: 4}
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedJson), w.Body.String())
}

func TestSubtractFromQuantity(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })
	api.Items[1] = api.Item{Name: "Book", Price: 12.99, Quantity: 5}
	requestDelta := api.Delta{QuantityDelta: -3}
	requestDeltaJson, _ := json.Marshal(requestDelta)

	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/item/1", strings.NewReader(string(requestDeltaJson)))
	router.ServeHTTP(w, req)

	expected := api.Item{Name: "Book", Price: 12.99, Quantity: 2}
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedJson), w.Body.String())
}

func TestSubtractUnderflow(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })
	api.Items[1] = api.Item{Name: "Book", Price: 12.99, Quantity: 5}
	requestDelta := api.Delta{QuantityDelta: -7}
	requestDeltaJson, _ := json.Marshal(requestDelta)

	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/item/1", strings.NewReader(string(requestDeltaJson)))
	router.ServeHTTP(w, req)

	expected := api.NewError("Quantity is too small to perform the operation.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 409, w.Code)
	assert.Equal(t, string(expectedJson), w.Body.String())
}

func TestAlterUnknownItem(t *testing.T) {
	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/item/1", strings.NewReader("{}"))
	router.ServeHTTP(w, req)

	expected := api.NewError("Item ID does not exist.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, string(expectedJson), w.Body.String())
}

func TestAlterInvalidId(t *testing.T) {
	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/item/invalid", strings.NewReader("{}"))
	router.ServeHTTP(w, req)

	expected := api.NewError("Invalid item ID.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, string(expectedJson), w.Body.String())
}
