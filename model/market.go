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
	SMA                 *float64 `json:"sma"`
	EMA                 *float64 `json:"ema"`
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

// SetIndicators ...
func (ms MarketPrices) SetIndicators() {
	ms.SetRSIs(14)
	ms.SetBollingerBands(20)
	ms.SetSMA(20)
	ms.SetEMA(20)
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

// SetSMA ...
func (ms MarketPrices) SetSMA(period int) {
	l := len(ms)
	prices := make([]float64, l)
	for i, m := range ms {
		prices[i] = m.Price
	}
	sma := talib.Sma(prices, period)
	for i, value := range sma {
		if value == 0 {
			continue
		}
		ms[i].SMA = util.ToPtr(value)
	}
}

// SetEMA ...
func (ms MarketPrices) SetEMA(period int) {
	l := len(ms)
	prices := make([]float64, l)
	for i, m := range ms {
		prices[i] = m.Price
	}
	ema := talib.Ema(prices, period)
	for i, value := range ema {
		if value == 0 {
			continue
		}
		ms[i].EMA = util.ToPtr(value)
	}
}

// ToPromptData ...
func (ms MarketPrices) ToPromptData() string {
	var result strings.Builder

	getFloat := func(p *float64) string {
		if p == nil {
			return "NaN"
		}
		return fmt.Sprintf("%f", *p)
	}

	for i, price := range ms {
		comma := ""
		if i != len(ms)-1 {
			comma = ","
		}
		result.WriteString(
			fmt.Sprintf("{\"date\": \"%s\", \"price\": %f, \"rsi\": %s, \"bollinger_band_upper\": %s, \"bollinger_band_middle\": %s, \"bollinger_band_lower\": %s, \"sma\": %s, \"ema\": %s}%s\n",
				price.DateUTC,
				price.Price,
				getFloat(price.RSI),
				getFloat(price.BollingerBandUpper),
				getFloat(price.BollingerBandMiddle),
				getFloat(price.BollingerBandLower),
				getFloat(price.SMA),
				getFloat(price.EMA),
				comma))
	}
	return result.String()
}
