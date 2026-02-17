package api

import (
	"inventory_manager/internal/db"
	"slices"
)

type Inventory struct {
	persistence db.InventoryPersistence
}

func (inventory Inventory) calculateValue(ids []uint) int64 {
	sum := int64(0)
	isEmptyQuery := ids == nil
	for id, item := range inventory.persistence.GetItems() {
		if isEmptyQuery || slices.Contains(ids, id) {
			sum += item.PriceCents * item.Quantity
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

func (inventory Inventory) alterQuantity(id uint, dQuantity int64) (*db.Item, *ErrorResponse) {
	var item, _ = inventory.persistence.GetItem(id)
	newQuantity := item.Quantity + dQuantity
	if newQuantity < 0 {
		return nil, new(NewError("Quantity is too small to perform the operation."))
	}
	item.Quantity = newQuantity
	inventory.persistence.SetItem(id, item)

	return &item, nil
}

func (inventory Inventory) itemExists(id uint) bool {
	_, exists := inventory.persistence.GetItem(id)
	return exists
}
