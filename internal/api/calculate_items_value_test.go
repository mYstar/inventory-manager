package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateItemsValue(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "GET", "/items/value?ids=1&ids=3", "")

	value := ValueResponse{}
	err := json.Unmarshal([]byte(response.Body.String()), &value)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, true, value.Success)
	assert.Equal(t, int64(330799), value.ValueCents)
}

func TestCalculateMissingIds(t *testing.T) {
	router := initTestRouterEmpty()
	response := sendTestRequest(router, "GET", "/items/value?ids=1&ids=2&ids=3", "")

	value := ValueResponse{}
	err := json.Unmarshal([]byte(response.Body.String()), &value)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, true, value.Success)
	assert.Equal(t, int64(0), value.ValueCents)
}

func TestCalculateInvalidIds(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "GET", "/items/value?ids=a&ids=2&ids=c", "")

	value := ErrorResponse{}
	err := json.Unmarshal([]byte(response.Body.String()), &value)

	assert.Nil(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, false, value.Success)
	assert.Equal(t, "One or more of given ids cannot be interpreted as uint.", value.Error)
}

func TestCalculateNegativeItemsValue(t *testing.T) {
	router := initTestRouter()
	sendTestRequest(router, "POST", "/item", `{"name": "Voucher", "price_cents": -1000, "quantity": 1}`)
	response := sendTestRequest(router, "GET", "/items/value?ids=1&ids=4", "")

	value := ValueResponse{}
	err := json.Unmarshal([]byte(response.Body.String()), &value)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, true, value.Success)
	assert.Equal(t, int64(99), value.ValueCents)
}
