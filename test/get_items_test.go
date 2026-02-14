package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetItems(t *testing.T) {
	originalItems := api.Items
	testItems := []api.Item{
		{ID: 1, Name: "Book", Price: 10.99, Quantity: 1},
		{ID: 2, Name: "Chair", Price: 24.89, Quantity: 4},
		{ID: 3, Name: "Laptop", Price: 1099.00, Quantity: 3},
	}
	api.Items = testItems
	t.Cleanup(func() { api.Items = originalItems })

	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items", nil)
	router.ServeHTTP(w, req)

	expected, _ := json.Marshal(testItems)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expected), w.Body.String())
}

func TestGetEmptyItems(t *testing.T) {

	router := api.SetupRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items", nil)
	router.ServeHTTP(w, req)

	expected := "[]"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, expected, w.Body.String())
}
