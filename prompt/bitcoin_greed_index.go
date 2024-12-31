package prompt

import (
	"magmar/model"
	"strings"
	"text/template"
)

type bitcoinGreedIndexPrompt struct {
	GreedIndex string
}

const bitcoinGreedIndexTemplate = `
	The Greed Index in Bitcoin, often referred to as the Crypto Fear & Greed Index, is a sentiment analysis tool designed to measure the emotional state of cryptocurrency market participants. 
	Extreme Fear is 0-24, Fear is 25-49, Neutral is 50, Greed is 50-74, Extreme Greed is 75-100
	{{.GreedIndex}}
`

// NewBitcoinGreedIndexPrompt ...
func NewBitcoinGreedIndexPrompt(greedIndex *model.GreedIndex) string {
	prompt := template.Must(template.New("bitcoin_orderbook").Parse(bitcoinOrderBookTemplate))

	var result strings.Builder
	prompt.Execute(&result, bitcoinGreedIndexPrompt{
		GreedIndex: greedIndex.ToPromptData(),
	})
	return result.String()
}
