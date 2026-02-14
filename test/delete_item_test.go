package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteItem(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })
	api.Items[1] = api.Item{Name: "Book", Price: 12.99, Quantity: 1}

	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/item/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestDeleteUnknownItem(t *testing.T) {
	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/item/1", nil)
	router.ServeHTTP(w, req)

	expected := api.NewError("Item ID does not exist.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, string(expectedJson), w.Body.String())
}

func TestDeleteInvalidId(t *testing.T) {
	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/item/invalid", nil)
	router.ServeHTTP(w, req)

	expected := api.NewError("Invalid item ID.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, string(expectedJson), w.Body.String())
}
