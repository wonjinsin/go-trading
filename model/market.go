package model

import (
	"fmt"
	"magmar/model/dao"
	"strings"
)

// MarketPrice ...
type MarketPrice struct {
	Market  string  `json:"market"`
	DateUTC string  `json:"date"`
	Price   float64 `json:"price"`
}

// MarketPrices ...
type MarketPrices []*MarketPrice

// NewMarketPriceByUpbit ...
func NewMarketPriceByUpbit(prices dao.UpbitMarketPrices) MarketPrices {
	var marketPrices MarketPrices
	for _, price := range prices {
		marketPrices = append(marketPrices, &MarketPrice{
			Market:  string(price.Market),
			DateUTC: price.CandleDateTimeUTC,
			Price:   price.TradePrice,
		})
	}
	return marketPrices
}

// ToPromptData ...
func (m MarketPrices) ToPromptData() string {
	var result strings.Builder
	for _, price := range m {
		result.WriteString(fmt.Sprintf("{\"date\": \"%s\", \"price\": %f}\n", price.DateUTC, price.Price))
	}
	return result.String()
}
