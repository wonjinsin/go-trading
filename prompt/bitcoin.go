package prompt

import (
	"html/template"
	"magmar/model"
	"strings"
)

type bitcoinPrompt struct {
	Strategy          string
	MarketData60Days  string
	MarketData24Hours string
	OrderBooks        string
	GreedIndex        string
	Newses            string
}

const bitcoinTemplate = `
    Your role is to serve as an advanced virtual assistant for Bitcoin trading, specifically for the KRW-BTC pair.
	Your objectives are to optimize profit margins and use a data-driven approach to guide trading decisions.
	I will ask every hour and trade by your decision.
	Please use the investment strategy and data to make a decision.
	
	Your investment strategy is based on the following principles:
	{{.Strategy}}

	Your Data is as follows:
	1. 60Days Market Datas
	{{.MarketData60Days}}

	2. 24hours Market Datas
	{{.MarketData24Hours}}

	3. Order Book Datas
	{{.OrderBooks}}

	4. Greed Index Data
	{{.GreedIndex}}

	5. News Datas
	{{.Newses}}
	
    Tell me decision whether to Buy, Sell, or Hold at the moment based on the provided data by json format.
	and percentage of how much percent should invest and reason why.
    Response example:
	{
		"decision": "Buy",
		"percent": 50,
		"reason": "After reviewing the current investment state and incorporating insights from both market analysis and recent crypto news, a bullish trend is evident. The EMA_10 has crossed above the SMA_10, a signal often associated with the initiation of an uptrend. This crossover, combined with our analysis of the current market sentiment being positively influenced by recent news articles, suggests increasing momentum and a strong opportunity for a profitable buy decision. This decision aligns with our risk management protocols, considering both the potential for profit and the current balance of the portfolio."
	},
	{
		"decision": "Buy",
		"percent": 55,
		"reason": "This decision to invest 25% of our portfolio in Bitcoin is predicated on a multi-faceted analysis incorporating market sentiment, technical indicators, and the prevailing crypto news landscape, alongside a prudent evaluation of our current investment state. A recent trend reversal has been identified, underscored by the EMA_10 decisively crossing above the SMA_10, heralding a bullish market phase. Concurrently, the RSI_14 reading has settled around 45, suggesting that Bitcoin is neither overbought nor oversold, thus presenting a compelling buy signal at the current price levels. Additionally, the Fear and Greed Index has recently dialed back from 'Extreme Greed' to 'Greed', signaling a cooling yet still positive investor sentiment, potentially pre-empting a market upswing. Notably, the latest crypto news analysis indicates a burgeoning confidence among institutional investors towards Bitcoin, particularly in light of regulatory clarity and advancements in blockchain technology, fostering a favorable environment for price appreciation. Furthermore, our portfolio's current allocation, with a balanced mix of BTC and KRW, coupled with an in-depth review of past trading decisions, suggests an opportune moment to augment our Bitcoin position. This strategic augmentation is designed to leverage anticipated market momentum while maintaining a vigilant stance on risk management, aiming to enhance our portfolio's profitability in alignment with our long-term investment objectives."
	}
    {
		"decision": "Sell",
		"percent": 60,
		"reason": "Upon detailed analysis of the asset's historical data and previous decision outcomes, it is evident that the asset is currently peaking near a historically significant resistance level. This observation is underscored by the RSI_14 indicator's ascent into overbought territory above 75, hinting at an overvaluation of the asset. Such overbought conditions are supported by a noticeable bearish divergence in the MACD, where despite the asset's price holding near its peak, the MACD line demonstrates a downward trajectory. This divergence aligns with a marked increase in trading volume, indicative of a potential buying climax which historically precedes market corrections. Reflecting on past predictions, similar conditions have often resulted in favorable sell outcomes, reinforcing the current decision to sell. Considering these factors - historical resistance alignment, overbought RSI_14, MACD bearish divergence, and peak trading volume - alongside a review of previous successful sell signals under comparable conditions, a strategic decision to sell 20% of the asset is recommended to leverage the anticipated market downturn and secure profits from the elevated price levels."
	},
	{
		"decision": "Hold",
		"percent": 0,
		"reason": "After a comprehensive review of the current market conditions, historical data, and previous decision outcomes, the present analysis indicates a complex trading environment. Although the MACD remains above its Signal Line, suggesting a potential buy signal, a notable decrease in the MACD Histogram's volume highlights diminishing momentum. This observation suggests caution, as weakening momentum could precede a market consolidation or reversal. Concurrently, the RSI_14 and SMA_10 indicators do not present clear signals, indicating a market in balance rather than one trending strongly in either direction. Furthermore, recent crypto news has introduced ambiguity into market sentiment, failing to provide a clear directional bias for the KRW-BTC pair. Considering these factors alongside a review of our portfolio's current state and in alignment with our risk management principles, the decision to hold reflects a strategic choice to preserve capital amidst market uncertainty. This cautious stance allows us to remain positioned for future opportunities while awaiting more definitive market signals."
	}
`

// NewBitcoinPrompt ...
func NewBitcoinPrompt(
	marketPrices model.MarketPrices,
	marketData24Hours model.MarketPrices,
	orderBooks model.OrderBooks,
	greedIndex *model.GreedIndex,
	newses model.Newses,
) string {
	prompt := template.Must(template.New("bitcoin").Parse(bitcoinTemplate))

	var result strings.Builder
	prompt.Execute(&result, bitcoinPrompt{
		Strategy:          NewBitcoinStrategyPrompt(),
		MarketData60Days:  NewBitcoinMarketData60DaysPrompt(marketPrices),
		MarketData24Hours: NewBitcoinMarketData24HoursPrompt(marketData24Hours),
		OrderBooks:        NewBitcoinOrderBookPrompt(orderBooks),
		GreedIndex:        NewBitcoinGreedIndexPrompt(greedIndex),
		Newses:            NewBitcoinNewsPrompt(newses),
	})
	return result.String()
}
