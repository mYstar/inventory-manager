package storage

import "errors"

// MemoryPersistence provides in-memory storage for items, facilitating operations on an `Items` collection.
// The data of this persistence provider is not persisted across restarts.
type MemoryPersistence struct {
	items Items
}

// NewMemoryPersistence returns a new MemoryPersistence instance.
func NewMemoryPersistence() *MemoryPersistence {
	return &MemoryPersistence{items: Items{}}
}

func (m *MemoryPersistence) GetItems() (Items, error) {
	return m.items, nil
}

func (m *MemoryPersistence) GetItem(id uint64) (Item, bool, error) {
	item, exists := m.items[id]
	return item, exists, nil
}

func (m *MemoryPersistence) AddItem(item Item) (uint64, error) {
	maxKey := uint64(0)
	for key := range m.items {
		if key > maxKey {
			maxKey = key
		}
	}
	newIdx := maxKey + 1
	m.items[newIdx] = item

	return newIdx, nil
}

func (m *MemoryPersistence) DeleteItem(id uint64) error {
	_, exists := m.items[id]
	if !exists {
		return errors.New("item id does not exist")
	}
	delete(m.items, id)
	return nil
}

func (m *MemoryPersistence) UpdateItem(id uint64, item Item) error {
	_, exists := m.items[id]
	if exists {
		m.items[id] = item
	}
	return nil
}
