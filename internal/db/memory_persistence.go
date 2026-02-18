package db

type MemoryPersistence struct {
	items Items
}

func NewMemoryPersistence() MemoryPersistence {
	return MemoryPersistence{items: Items{}}
}

func (memory MemoryPersistence) GetItems() (Items, error) {
	return memory.items, nil
}

func (memory MemoryPersistence) GetItem(id uint) (Item, bool, error) {
	item, exists := memory.items[id]
	return item, exists, nil
}

func (memory MemoryPersistence) AddItem(item Item) (uint, error) {
	maxKey := uint(0)
	for key := range memory.items {
		if key > maxKey {
			maxKey = key
		}
	}
	newIdx := maxKey + 1
	memory.items[newIdx] = item

	return newIdx, nil
}

func (memory MemoryPersistence) DeleteItem(id uint) error {
	delete(memory.items, id)
	return nil
}

func (memory MemoryPersistence) SetItem(id uint, item Item) error {
	memory.items[id] = item
	return nil
}
