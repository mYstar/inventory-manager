package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetItems(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "GET", "/items", "")

	// TODO fix this test
	//	expected, _ := json.Marshal(api.Inventory)
	assert.Equal(t, 200, response.Code)
	//	assert.Equal(t, string(expected), response.Body.String())
}

func TestGetEmptyItems(t *testing.T) {
	router := initTestRouterEmpty()
	response := sendTestRequest(router, "GET", "/items", "")

	expected := "{}"
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, expected, response.Body.String())
}
