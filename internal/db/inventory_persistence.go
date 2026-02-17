package db

type Item struct {
	Name       string `json:"name"`
	PriceCents int64  `json:"price_cents"`
	Quantity   int64  `json:"quantity"`
}

type Items map[uint]Item

type InventoryPersistence interface {
	GetItems() Items
	GetItem(uint) (Item, bool)
	AddItem(item Item) uint
	SetItem(uint, Item)
	DeleteItem(uint)
}
