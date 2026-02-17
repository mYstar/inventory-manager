package db

type MemoryPersistence struct {
	items Items
}

func NewMemoryPersistence() MemoryPersistence {
	return MemoryPersistence{items: Items{}}
}

func (memory MemoryPersistence) GetItems() Items {
	return memory.items
}

func (memory MemoryPersistence) GetItem(id uint) (Item, bool) {
	item, exists := memory.items[id]
	return item, exists
}

func (memory MemoryPersistence) AddItem(item Item) uint {
	maxKey := uint(0)
	for key := range memory.items {
		if key > maxKey {
			maxKey = key
		}
	}
	newIdx := maxKey + 1
	memory.items[newIdx] = item

	return newIdx
}

func (memory MemoryPersistence) DeleteItem(id uint) {
	delete(memory.items, id)
}

func (memory MemoryPersistence) SetItem(id uint, item Item) {
	memory.items[id] = item
}
