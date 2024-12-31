package prompt

import (
	"magmar/model"
	"strings"
	"text/template"
)

type bitcoinMarketData24HoursPrompt struct {
	MarketPrices string
}

const bitcoinMarketData24HoursTemplate = `
	MarketPrice for 24 hours per an hour
	I will offer date, price, rsi of range 14, bollinger_band_upper of range 20, bollinger_band_middle of range 20, bollinger_band_lower of range 20, sma of range 20, ema of range 20
	Value is NaN if not available
	{{.MarketPrices}}
`

// NewBitcoinMarketData24HoursPrompt ...
func NewBitcoinMarketData24HoursPrompt(marketPrices model.MarketPrices) string {
	prompt := template.Must(template.New("bitcoin_market_data_24hours").Parse(bitcoinMarketData24HoursTemplate))

	var result strings.Builder
	prompt.Execute(&result, bitcoinMarketData24HoursPrompt{
		MarketPrices: marketPrices.ToPromptData(),
	})
	return result.String()
}
