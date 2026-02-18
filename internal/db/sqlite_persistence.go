package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type SqlitePersistence struct {
	db *sql.DB
}

func NewSqlitePersistence(dbFile string) (*SqlitePersistence, error) {

	db, err := sql.Open("sqlite", "data/"+dbFile)
	if err != nil {
		return nil, err
	}

	// This ensures the table exists every time the app starts
	initQuery := `
	CREATE TABLE IF NOT EXISTS inventory (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(128) NOT NULL,
	    quantity INTEGER NOT NULL,
	    price_cents INTEGER NOT NULL
	);`
	_, err = db.Exec(initQuery)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return new(SqlitePersistence{db}), nil
}

func (memory SqlitePersistence) GetItems() (Items, error) {
	getAllQuery := "SELECT * FROM inventory"
	rows, err := memory.db.Query(getAllQuery)
	if err != nil {
		return Items{}, err
	}
	defer rows.Close()

	items := Items{}
	for rows.Next() {
		id := uint(0)
		item := Item{}
		err = rows.Scan(&id, &item.Name, &item.Quantity, &item.PriceCents)
		if err != nil {
			return Items{}, err
		}
		items[id] = item
	}

	return items, nil
}

func (memory SqlitePersistence) GetItem(id uint) (Item, bool, error) {
	getQuery := "SELECT name, quantity, price_cents FROM inventory WHERE id = ?"
	rows, err := memory.db.Query(getQuery, id)
	if err != nil {
		return Item{}, false, err
	}
	defer rows.Close()

	exists := rows.Next()
	if !exists {
		return Item{}, false, nil
	}
	var item Item
	err = rows.Scan(&item.Name, &item.Quantity, &item.PriceCents)
	if err != nil {
		return Item{}, false, err
	}
	return item, true, nil
}

func (memory SqlitePersistence) AddItem(item Item) (uint, error) {
	setQuery := "INSERT INTO inventory (name, quantity, price_cents) VALUES (?, ?, ?)"
	result, err := memory.db.Exec(setQuery, item.Name, item.Quantity, item.PriceCents)
	if err != nil {
		return 0, err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(newId), nil
}

func (memory SqlitePersistence) DeleteItem(id uint) error {
	deleteQuery := "DELETE FROM inventory WHERE id = ?"
	_, err := memory.db.Exec(deleteQuery, id)

	return err
}

func (memory SqlitePersistence) SetItem(id uint, item Item) error {
	setQuery := "UPDATE inventory SET name = ?, quantity = ?, price_cents = ? WHERE id = ?"
	_, err := memory.db.Exec(setQuery, item.Name, item.Quantity, item.PriceCents, id)

	return err
}
