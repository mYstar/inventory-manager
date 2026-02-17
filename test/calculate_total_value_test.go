package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTotalValue(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "GET", "/items/value", "")

	value := api.ValueResponse{}
	err := json.Unmarshal([]byte(response.Body.String()), &value)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, true, value.Success)
	assert.Equal(t, float32(3407.55), value.Value)
}

func TestCalculateTotalValueEmpty(t *testing.T) {
	router := initTestRouterEmpty()
	response := sendTestRequest(router, "GET", "/items/value", "")

	value := api.ValueResponse{}
	err := json.Unmarshal([]byte(response.Body.String()), &value)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, true, value.Success)
	assert.Equal(t, float32(0.0), value.Value)
}
