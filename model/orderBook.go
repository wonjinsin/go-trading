package model

// OrderBook ...
type OrderBook struct {
	AskPrice uint64  `json:"ask_price"`
	AskSize  uint64  `json:"ask_size"`
	BidPrice float64 `json:"bid_price"`
	BidSize  float64 `json:"bid_size"`
}
