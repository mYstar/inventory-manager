package api

type DeltaRequest struct {
	QuantityDelta *int64 `json:"quantity_delta"`
}

type IdsQuery struct {
	Ids []uint `form:"ids"`
}

type ItemResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	PriceCents int64  `json:"price_cents"`
	Quantity   int64  `json:"quantity"`
}

type ValueResponse struct {
	Success    bool  `json:"success"`
	ValueCents int64 `json:"value_cents"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
