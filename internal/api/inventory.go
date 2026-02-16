package api

import (
	"slices"
)

type Item struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity uint    `json:"quantity"`
}

type Items map[uint]Item

var Inventory = Items{}

func (inventory Items) calculateValue(ids []uint) float32 {
	sum := float32(0.0)
	isEmptyQuery := ids == nil
	for id, item := range inventory {
		if isEmptyQuery || slices.Contains(ids, id) {
			sum += item.Price * float32(item.Quantity)
		}
	}

	return sum
}

func (inventory Items) createItem(item Item) (uint, Item) {
	maxKey := uint(0)
	for key := range inventory {
		if key > maxKey {
			maxKey = key
		}
	}
	newIdx := maxKey + 1
	inventory[newIdx] = item

	return newIdx, inventory[newIdx]
}

func (inventory Items) delete(id uint) {
	delete(inventory, id)
}

func (inventory Items) alterQuantity(id uint, dQuantity int) (*Item, *ErrorResponse) {
	var item = Inventory[id]
	switch {
	case dQuantity < 0:
		if uint(-dQuantity) > item.Quantity {

			return nil, new(NewError("Quantity is too small to perform the operation."))
		}
		item.Quantity -= uint(-dQuantity)
	default:
		item.Quantity += uint(dQuantity)
	}
	Inventory[id] = item

	return &item, nil
}

func (inventory Items) itemExists(id uint) bool {
	_, exists := Inventory[id]

	return exists
}
