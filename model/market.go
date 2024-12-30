package model

import (
	"fmt"
	"magmar/model/dao"
	"magmar/util"
	"strings"

	"github.com/markcheno/go-talib"
)

// MarketPrice ...
type MarketPrice struct {
	Market              string   `json:"market"`
	DateUTC             string   `json:"date"`
	Price               float64  `json:"price"`
	RSI                 *float64 `json:"rsi"`
	BollingerBandUpper  *float64 `json:"bollinger_band_upper"`
	BollingerBandMiddle *float64 `json:"bollinger_band_middle"`
	BollingerBandLower  *float64 `json:"bollinger_band_lower"`
}

// MarketPrices ...
type MarketPrices []*MarketPrice

// NewMarketPriceByUpbit ...
func NewMarketPriceByUpbit(prices dao.UpbitMarketPrices) MarketPrices {
	var marketPrices MarketPrices
	// order by date
	for i := len(prices) - 1; i >= 0; i-- {
		price := prices[i]
		marketPrices = append(marketPrices, &MarketPrice{
			Market:  string(price.Market),
			DateUTC: price.CandleDateTimeUTC,
			Price:   price.TradePrice,
		})
	}
	return marketPrices
}

// SetRSIs ...
func (ms MarketPrices) SetRSIs(period int) {
	l := len(ms)
	prices := make([]float64, l)
	for i, m := range ms {
		prices[i] = m.Price
	}
	rsis := talib.Rsi(prices, period)
	for i, value := range rsis {
		if value == 0 {
			continue
		}
		ms[i].RSI = util.ToPtr(value)
	}
}

// SetBollingerBands ...
func (ms MarketPrices) SetBollingerBands(period int) {
	l := len(ms)
	prices := make([]float64, l)
	for i, m := range ms {
		prices[i] = m.Price
	}
	up, mid, low := talib.BBands(prices, period, 2, 2, talib.SMA)
	for i := range prices {
		if up[i] != 0 {
			ms[i].BollingerBandUpper = util.ToPtr(up[i])
		}
		if mid[i] != 0 {
			ms[i].BollingerBandMiddle = util.ToPtr(mid[i])
		}
		if low[i] != 0 {
			ms[i].BollingerBandLower = util.ToPtr(low[i])
		}
	}
}

// ToPromptData ...
func (ms MarketPrices) ToPromptData() string {
	var result strings.Builder
	for _, price := range ms {
		result.WriteString(fmt.Sprintf("{\"date\": \"%s\", \"price\": %f}\n", price.DateUTC, price.Price))
	}
	return result.String()
}
