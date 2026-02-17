package test

import (
	"encoding/json"
	"inventory_manager/internal/api"
	"inventory_manager/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateItem(t *testing.T) {
	router := initTestRouterEmpty()
	testItem := db.Item{
		Name:     "Book",
		Price:    10.99,
		Quantity: 1,
	}
	testItemJson, _ := json.Marshal(testItem)

	response := sendTestRequest(router, "POST", "/item", string(testItemJson))

	createdItemJson := response.Body.String()
	createdItem := api.ItemResponse{}
	err := json.Unmarshal([]byte(createdItemJson), &createdItem)

	assert.Nil(t, err)
	assert.Equal(t, 201, response.Code)
	assert.Equal(t, uint(1), createdItem.ID)
	assert.Equal(t, testItem.Name, createdItem.Name)
	assert.Equal(t, testItem.Quantity, createdItem.Quantity)
	assert.Equal(t, testItem.Price, createdItem.Price)
}

func TestCreateItemExistingValues(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "DELETE", "/item/2", "")
	testItem := db.Item{
		Name:     "Sofa",
		Price:    509.99,
		Quantity: 3,
	}
	testItemJson, _ := json.Marshal(testItem)

	// check if creating works
	response = sendTestRequest(router, "POST", "/item", string(testItemJson))

	createdItemJson := response.Body.String()
	createdItem := api.ItemResponse{}
	err := json.Unmarshal([]byte(createdItemJson), &createdItem)

	assert.Nil(t, err)
	assert.Equal(t, 201, response.Code)
	assert.Equal(t, uint(4), createdItem.ID)
	assert.Equal(t, testItem.Name, createdItem.Name)
	assert.Equal(t, testItem.Quantity, createdItem.Quantity)
	assert.Equal(t, testItem.Price, createdItem.Price)

	// check if the new item is in the inventory
	response = sendTestRequest(router, "GET", "/items", "")

	// TODO fix this test
	//initialItems := maps.NewMap[uint, db.Item]()
	//initialItems[4] = testItem
	//expected, _ := json.Marshal(initialItems)
	//assert.Equal(t, 200, response.Code)
	//assert.Equal(t, string(expected), response.Body.String())
}

func TestCreateItemWithInvalidJson(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "POST", "/item", "{")

	expected := api.NewError("Request body is not valid JSON.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}

func TestCreateItemWithMissingValues(t *testing.T) {
	router := initTestRouter()
	response := sendTestRequest(router, "POST", "/item", "{\"name\": \"Book\", \"price\": 10.99}")

	expected := api.NewError("Missing item data.")
	expectedJson, _ := json.Marshal(expected)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, string(expectedJson), response.Body.String())
}
