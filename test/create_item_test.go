package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateItem(t *testing.T) {
	originalItems := api.Items
	t.Cleanup(func() { api.Items = originalItems })
	testItem := struct {
		Name     string  `json:"name"`
		Price    float32 `json:"price"`
		Quantity uint    `json:"quantity"`
	}{
		Name:     "Book",
		Price:    10.99,
		Quantity: 1,
	}
	testItemJson, _ := json.Marshal(testItem)

	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/item", strings.NewReader(string(testItemJson)))
	router.ServeHTTP(w, req)

	createdItemJson := w.Body.String()
	createdItem := api.Item{}
	err := json.Unmarshal([]byte(createdItemJson), &createdItem)

	assert.Nil(t, err)
	assert.Equal(t, 201, w.Code)
	assert.NotEqual(t, 0, createdItem.ID)
	assert.Equal(t, testItem.Name, createdItem.Name)
	assert.Equal(t, testItem.Quantity, createdItem.Quantity)
	assert.Equal(t, testItem.Price, createdItem.Price)
}

func TestCreateItemWithInvalidJson(t *testing.T) {
	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/item", strings.NewReader("{"))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestCreateItemWithMissingValues(t *testing.T) {
	router := api.SetupRoutes()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/item", strings.NewReader("{\"name\": \"Book\", \"price\": 10.99}"))
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
