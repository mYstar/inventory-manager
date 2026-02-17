package test

import (
	"inventory_manager/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetItems(t *testing.T) {
	router := initTestRouter()
	expected := db.Items{
		1: {Name: "Book", Price: 10.99, Quantity: 1},
		2: {Name: "Chair", Price: 24.89, Quantity: 4},
		3: {Name: "Laptop", Price: 1099.00, Quantity: 3},
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
