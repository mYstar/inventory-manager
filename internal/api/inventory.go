package api

type Item struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity uint    `json:"quantity"`
}

var Items = make(map[uint]Item)
