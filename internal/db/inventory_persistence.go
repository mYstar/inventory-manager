package db

type Item struct {
	Name       string `json:"name"`
	PriceCents int64  `json:"price_cents"`
	Quantity   int64  `json:"quantity"`
}

type Items map[uint]Item

type InventoryPersistence interface {
	GetItems() (Items, error)
	GetItem(uint) (Item, bool, error)
	AddItem(item Item) (uint, error)
	SetItem(uint, Item) error
	DeleteItem(uint) error
}
