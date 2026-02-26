package storage

import (
	"database/sql"
	"errors"
	"strings"

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
	    quantity INTEGER NOT NULL CHECK(quantity >= 0),
	    price_cents INTEGER NOT NULL
	);`
	_, err = db.Exec(initQuery)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return &SqlitePersistence{db}, nil
}

func (s *SqlitePersistence) GetItems() (Items, error) {
	getAllQuery := "SELECT * FROM inventory"
	rows, err := s.db.Query(getAllQuery)
	if err != nil {
		return Items{}, err
	}
	defer rows.Close()

	items := Items{}
	for rows.Next() {
		id := uint64(0)
		item := Item{}
		err = rows.Scan(&id, &item.Name, &item.Quantity, &item.PriceCents)
		if err != nil {
			return Items{}, err
		}
		items[id] = item
	}

	return items, nil
}

func (s *SqlitePersistence) GetItem(id uint64) (Item, bool, error) {
	getQuery := "SELECT name, quantity, price_cents FROM inventory WHERE id = ?"
	row := s.db.QueryRow(getQuery, id)

	var item Item
	err := row.Scan(&item.Name, &item.Quantity, &item.PriceCents)
	if err != nil {
		return Item{}, false, err
	}
	return item, true, nil
}

func (s *SqlitePersistence) AddItem(item Item) (uint64, Item, error) {
	setQuery := "INSERT INTO inventory (name, quantity, price_cents) VALUES (?, ?, ?)"
	result, err := s.db.Exec(setQuery, item.Name, item.Quantity, item.PriceCents)
	if err != nil {
		return 0, Item{}, err
	}

	newId, err := result.LastInsertId()
	if err != nil {
		return 0, Item{}, err
	}
	return uint64(newId), item, nil
}

func (s *SqlitePersistence) DeleteItem(id uint64) error {
	deleteQuery := "DELETE FROM inventory WHERE id = ?"
	_, err := s.db.Exec(deleteQuery, id)

	return err
}

func (s *SqlitePersistence) AlterQuantityBy(id uint64, deltaQuantity int64) (Item, error) {
	setQuery := "UPDATE inventory SET quantity = quantity + ? WHERE id = ? RETURNING name, quantity, price_cents"
	row := s.db.QueryRow(setQuery, deltaQuantity, id)

	var item Item
	err := row.Scan(&item.Name, &item.Quantity, &item.PriceCents)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return item, errors.New("item id does not exist")
		}
		if strings.HasPrefix(err.Error(), "constraint failed: CHECK constraint failed: quantity") {
			return Item{}, errors.New("quantity is too small to perform the operation")
		}
		return Item{}, err
	}

	return item, nil
}
