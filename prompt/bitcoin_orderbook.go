package prompt

import (
	"magmar/model"
	"strings"
	"text/template"
)

type bitcoinOrderBookPrompt struct {
	OrderBook string
}

const bitcoinOrderBookTemplate = `
	The order book represents the real-time list of buy and sell orders in the market.
	t provides a snapshot of market activity and liquidity, showing traders the supply and demand dynamics for an asset.
	I will offer ask_price, ask_size, bid_price, bid_size
	{{.OrderBook}}
`

// NewBitcoinOrderBookPrompt ...
func NewBitcoinOrderBookPrompt(orderBooks model.OrderBooks) string {
	prompt := template.Must(template.New("bitcoin_orderbook").Parse(bitcoinOrderBookTemplate))

	var result strings.Builder
	prompt.Execute(&result, bitcoinOrderBookPrompt{
		OrderBook: orderBooks.ToPromptData(),
	})
	return result.String()
}
