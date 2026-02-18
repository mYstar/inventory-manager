package api

import (
	"inventory_manager/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetItems(t *testing.T) {
	router := initTestRouter()
	expected := storage.Items{
		1: {Name: "Book", PriceCents: 1099, Quantity: 1},
		2: {Name: "Chair", PriceCents: 2489, Quantity: 4},
		3: {Name: "Laptop", PriceCents: 109900, Quantity: 3},
	}

	assertInventoryEquals(t, router, expected)
}

func TestGetEmptyItems(t *testing.T) {
	router := initTestRouterEmpty()
	response := sendTestRequest(router, "GET", "/items", "")

	expected := "{}"
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, expected, response.Body.String())
}
