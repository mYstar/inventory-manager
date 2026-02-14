package api

type Item struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity uint    `json:"quantity"`
}

var Items = make([]Item, 0)
