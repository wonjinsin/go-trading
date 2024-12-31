package model

import (
	"fmt"
	"strings"
)

// OrderBook ...
type OrderBook struct {
	AskPrice uint64  `json:"ask_price"`
	AskSize  uint64  `json:"ask_size"`
	BidPrice float64 `json:"bid_price"`
	BidSize  float64 `json:"bid_size"`
}

// OrderBooks ...
type OrderBooks []*OrderBook

// ToPromptData ...
func (os OrderBooks) ToPromptData() string {
	var result strings.Builder
	for i, o := range os {
		comma := ""
		if i != len(os)-1 {
			comma = ","
		}
		result.WriteString(fmt.Sprintf("{\"ask_price\": %d, \"ask_size\": %d, \"bid_price\": %f, \"bid_size\": %f}%s\n", o.AskPrice, o.AskSize, o.BidPrice, o.BidSize, comma))
	}
	return result.String()
}
