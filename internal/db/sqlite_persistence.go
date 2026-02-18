package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type SqlitePersistence struct {
	db *sql.DB
}

func NewSqlitePersistence() SqlitePersistence {

	db, err := sql.Open("sqlite", "data/default.db")
	if err != nil {
		panic(err)
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
		panic(err)
	}

	return SqlitePersistence{db}
}

func (memory SqlitePersistence) GetItems() Items {
	getAllQuery := "SELECT * FROM inventory"
	rows, err := memory.db.Query(getAllQuery)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	items := Items{}
	for rows.Next() {
		id := uint(0)
		item := Item{}
		err = rows.Scan(&id, &item.Name, &item.Quantity, &item.PriceCents)
		if err != nil {
			panic(err)
		}
		items[id] = item
	}

	return items
}

func (memory SqlitePersistence) GetItem(id uint) (Item, bool) {
	getQuery := "SELECT name, quantity, price_cents FROM inventory WHERE id = ?"
	rows, err := memory.db.Query(getQuery, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	rows.Next()
	var item Item
	err = rows.Scan(&item.Name, &item.Quantity, &item.PriceCents)
	if err != nil {
		panic(err)
	}
	return item, true
}

func (memory SqlitePersistence) AddItem(item Item) uint {
	setQuery := "INSERT INTO inventory (name, quantity, price_cents) VALUES (?, ?, ?)"
	result, err := memory.db.Exec(setQuery, item.Name, item.Quantity, item.PriceCents)
	if err != nil {
		panic(err)
	}

	newId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	return uint(newId)
}

func (memory SqlitePersistence) DeleteItem(id uint) {
	deleteQuery := "DELETE FROM inventory WHERE id = ?"
	_, err := memory.db.Exec(deleteQuery, id)
	if err != nil {
		panic(err)
	}
}

func (memory SqlitePersistence) SetItem(id uint, item Item) {
	setQuery := "UPDATE inventory SET name = ?, quantity = ?, price_cents = ? WHERE id = ?"
	_, err := memory.db.Exec(setQuery, item.Name, item.Quantity, item.PriceCents, id)
	if err != nil {
		panic(err)
	}
}
