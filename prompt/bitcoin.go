package prompt

import (
	"magmar/model"
	"strings"
	"text/template"
)

type bitcoinPrompt struct {
	MarketData string
}

const bitcoinTemplate = `
    You're a Bitcoin expert.,
    Tell me whether to buy, sell, or hold at the moment based on the provided chart data,
    Response example:,
    {"decision": "Buy", "reason": "Some technical reason"},
    {"decision": "Sell", "reason": "Some technical reason"},
    {"decision": "Hold", "reason": "Some technical reason"},
    Current Bitcoin Data:
	{{.MarketData}}
`

// NewBitcoinPrompt ...
func NewBitcoinPrompt(marketPrices model.MarketPrices) string {
	prompt := template.Must(template.New("bitcoin").Parse(bitcoinTemplate))

	var result strings.Builder
	prompt.Execute(&result, bitcoinPrompt{MarketData: marketPrices.ToPromptData()})
	return result.String()
}
