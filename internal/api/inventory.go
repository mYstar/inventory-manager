package api

import (
	"errors"
	"inventory_manager/internal/storage"
	"slices"
)

// Inventory represents the inventory system that manages items using the provided persistence layer.
type Inventory struct {
	persistence storage.InventoryPersistence
}

// NewInventory returns a new Inventory instance using the provided persistence layer.
func NewInventory(persistence storage.InventoryPersistence) *Inventory {
	return new(Inventory{persistence: persistence})
}

// GetItems retrieves all items from the inventory. Returns the items or an error if retrieving fails.
func (inventory *Inventory) GetItems() (storage.Items, error) {
	return inventory.persistence.GetItems()
}

// calculateValue calculates the value of the specified items or all items if the IDs list is nil.
// Returns the total value in cents or 0.0 and an error if retrieving items fails.
func (inventory *Inventory) calculateValue(ids []uint64) (int64, error) {
	sum := int64(0)
	isEmptyQuery := ids == nil
	items, err := inventory.persistence.GetItems()
	if err != nil {
		return 0, err
	}

	for id, item := range items {
		if isEmptyQuery || slices.Contains(ids, id) {
			sum += item.PriceCents * item.Quantity
		}
	}

	return sum, err
}

// createItem adds a new item to the inventory.
// Returns ID, the new item, or an error if adding or retrieving fails.
func (inventory *Inventory) createItem(item storage.Item) (uint64, storage.Item, error) {
	newIdx, err := inventory.persistence.AddItem(item)
	if err != nil {
		return 0, storage.Item{}, err
	}
	newItem, exists, err := inventory.persistence.GetItem(newIdx)
	if !exists {
		return 0, storage.Item{}, errors.New("item not found after creation")
	}
	return newIdx, newItem, err
}

// delete removes an item from the inventory by its ID.
// Returns an error if the operation fails.
func (inventory *Inventory) delete(id uint64) error {
	return inventory.persistence.DeleteItem(id)
}

// alterQuantity adjusts the quantity of an item in the inventory by the specified delta.
// Returns the altered item or an error if the operation is invalid or fails.
func (inventory *Inventory) alterQuantity(id uint64, dQuantity int64) (storage.Item, error) {
	item, exists, err := inventory.persistence.GetItem(id)
	if !exists {
		return storage.Item{}, errors.New("item id does not exist")
	}
	if err != nil {
		return storage.Item{}, err
	}
	newQuantity := item.Quantity + dQuantity
	if newQuantity < 0 {
		return storage.Item{}, errors.New("quantity is too small to perform the operation")
	}
	item.Quantity = newQuantity
	err = inventory.persistence.UpdateItem(id, item)

	return item, err
}
