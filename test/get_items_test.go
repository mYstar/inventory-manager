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

func TestGetItems(t *testing.T) {
	originalItems := maps.Clone(api.Items)
	t.Cleanup(func() { api.Items = originalItems })

	testItems := make(map[uint]api.Item)
	testItems[1] = api.Item{Name: "Book", Price: 10.99, Quantity: 1}
	testItems[2] = api.Item{Name: "Chair", Price: 24.89, Quantity: 4}
	testItems[3] = api.Item{Name: "Laptop", Price: 1099.00, Quantity: 3}
	api.Items = testItems

	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items", nil)
	router.ServeHTTP(w, req)

	expected, _ := json.Marshal(api.Items)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expected), w.Body.String())
}

func TestGetEmptyItems(t *testing.T) {

	router := api.SetupRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items", nil)
	router.ServeHTTP(w, req)

	expected := "{}"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, expected, w.Body.String())
}
