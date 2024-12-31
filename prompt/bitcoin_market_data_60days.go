package prompt

import (
	"magmar/model"
	"strings"
	"text/template"
)

type bitcoinMarketData60DaysPrompt struct {
	MarketPrices string
}

const bitcoinMarketData60DaysTemplate = `
	MarketPrice for 60 days per a day
	I will offer date, price, rsi of range 14, bollinger_band_upper of range 20, bollinger_band_middle of range 20, bollinger_band_lower of range 20, sma of range 20, ema of range 20
	Value is NaN if not available
	{{.MarketPrices}}
`

// NewBitcoinMarketData60DaysPrompt ...
func NewBitcoinMarketData60DaysPrompt(marketPrices model.MarketPrices) string {
	prompt := template.Must(template.New("bitcoin_market_data_60days").Parse(bitcoinMarketData60DaysTemplate))

	var result strings.Builder
	prompt.Execute(&result, bitcoinMarketData60DaysPrompt{
		MarketPrices: marketPrices.ToPromptData(),
	})
	return result.String()
}
