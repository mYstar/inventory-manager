package api

// DeltaRequest represents a request to change the quantity of an item.
type DeltaRequest struct {
	QuantityDelta *int64 `json:"quantity_delta"`
}

// IdsQuery represents a URL query that identifies a list of items by their IDs.
type IdsQuery struct {
	Ids []uint `form:"ids"`
}

// ItemResponse represents a response containing the ID and all data of an item.
type ItemResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	PriceCents int64  `json:"price_cents"`
	Quantity   int64  `json:"quantity"`
}

// ValueResponse represents a response containing the total value of a list of items.
type ValueResponse struct {
	Success    bool  `json:"success"`
	ValueCents int64 `json:"value_cents"`
}

// ErrorResponse represents a response containing an error message.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
