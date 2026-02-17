package api

import (
	"inventory_manager/internal/db"
	"slices"
)

type Inventory struct {
	persistence db.InventoryPersistence
}

func (inventory Inventory) calculateValue(ids []uint) float32 {
	sum := float32(0.0)
	isEmptyQuery := ids == nil
	for id, item := range inventory.persistence.GetItems() {
		if isEmptyQuery || slices.Contains(ids, id) {
			sum += item.Price * float32(item.Quantity)
		}
	}

	return sum
}

func (inventory Inventory) createItem(item db.Item) (uint, db.Item) {
	newIdx := inventory.persistence.AddItem(item)
	newItem, _ := inventory.persistence.GetItem(newIdx)
	return newIdx, newItem
}

func (inventory Inventory) delete(id uint) {
	inventory.persistence.DeleteItem(id)
}

func (inventory Inventory) alterQuantity(id uint, dQuantity int) (*db.Item, *ErrorResponse) {
	var item, _ = inventory.persistence.GetItem(id)
	switch {
	case dQuantity < 0:
		if uint(-dQuantity) > item.Quantity {

			return nil, new(NewError("Quantity is too small to perform the operation."))
		}
		item.Quantity -= uint(-dQuantity)
	default:
		item.Quantity += uint(dQuantity)
	}
	inventory.persistence.SetItem(id, item)

	return &item, nil
}

func (inventory Inventory) itemExists(id uint) bool {
	_, exists := inventory.persistence.GetItem(id)
	return exists
}
