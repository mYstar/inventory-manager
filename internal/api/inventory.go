package api

import (
	"errors"
	"inventory_manager/internal/db"
	"slices"
)

type Inventory struct {
	persistence db.InventoryPersistence
}

func NewInventory(persistence db.InventoryPersistence) Inventory {
	return Inventory{persistence: persistence}
}

func (inventory Inventory) calculateValue(ids []uint) (int64, error) {
	sum := int64(0)
	isEmptyQuery := ids == nil
	items, err := inventory.persistence.GetItems()

	for id, item := range items {
		if isEmptyQuery || slices.Contains(ids, id) {
			sum += item.PriceCents * item.Quantity
		}
	}

	return sum, err
}

func (inventory Inventory) createItem(item db.Item) (uint, db.Item, error) {
	newIdx, err := inventory.persistence.AddItem(item)
	if err != nil {
		return 0, db.Item{}, err
	}
	newItem, exists, err := inventory.persistence.GetItem(newIdx)
	if !exists {
		return 0, db.Item{}, errors.New("item not found after creation")
	}
	return newIdx, newItem, err
}

func (inventory Inventory) delete(id uint) error {
	return inventory.persistence.DeleteItem(id)
}

func (inventory Inventory) alterQuantity(id uint, dQuantity int64) (db.Item, error) {
	item, exists, err := inventory.persistence.GetItem(id)
	if !exists {
		return db.Item{}, errors.New("item id does not exist")
	}
	if err != nil {
		return db.Item{}, err
	}
	newQuantity := item.Quantity + dQuantity
	if newQuantity < 0 {
		return db.Item{}, errors.New("quantity is too small to perform the operation")
	}
	item.Quantity = newQuantity
	err = inventory.persistence.SetItem(id, item)

	return item, err
}

func (inventory Inventory) itemExists(id uint) bool {
	_, exists, _ := inventory.persistence.GetItem(id)
	return exists
}
