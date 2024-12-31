package prompt

import (
	"magmar/model"
	"strings"
	"text/template"
)

type bitcoinNewsPrompt struct {
	News string
}

const bitcoinNewsTemplate = `
	Recent 5 news about Bitcoin and prediction of future market
	{{.News}}
`

// NewBitcoinNewsPrompt ...
func NewBitcoinNewsPrompt(newses model.Newses) string {
	prompt := template.Must(template.New("bitcoin_news").Parse(bitcoinNewsTemplate))

	var result strings.Builder
	prompt.Execute(&result, bitcoinNewsPrompt{
		News: newses.ToPromptData(),
	})
	return result.String()
}
