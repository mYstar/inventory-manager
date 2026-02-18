package storage

// Item represents a single item in the inventory.
type Item struct {
	Name       string `json:"name"`
	PriceCents int64  `json:"price_cents"`
	Quantity   int64  `json:"quantity"`
}

// Items represents a map of items indexed by their IDs.
type Items map[uint]Item

// InventoryPersistence defines the interface for managing inventory operations.
type InventoryPersistence interface {
	// GetItems returns all items in the inventory or an error if retrieving fails.
	GetItems() (Items, error)
	// GetItem returns an item with the specified ID or an error if retrieving fails.
	GetItem(uint) (Item, bool, error)
	// AddItem adds a new item to the inventory and returns its ID or an error if adding fails.
	AddItem(item Item) (uint, error)
	// UpdateItem updates an existing item in the inventory or returns an error if updating fails.
	// Has no effect if the item does not exist.
	UpdateItem(uint, Item) error
	// DeleteItem removes an item from the inventory or returns an error if deleting fails.
	DeleteItem(uint) error
}
