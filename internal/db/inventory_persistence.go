package db

type Item struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity uint    `json:"quantity"`
}

type Items map[uint]Item

type InventoryPersistence interface {
	GetItems() Items
	GetItem(uint) (Item, bool)
	AddItem(item Item) uint
	SetItem(uint, Item)
	DeleteItem(uint)
}
