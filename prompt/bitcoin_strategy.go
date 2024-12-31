package prompt

// BitcoinStrategyTemplate ...
const BitcoinStrategyTemplate = `
1. Basic Concept of Bollinger Bands
- 20-day moving average as center line
- Upper and Lower bands calculated as 2 times the standard deviation
- 95% probability that price stays within the bands

Adjust band range by changing standard deviation:
- Std dev 1 → 68% probability
- Std dev 3 → 99% probability

2. Basic Trading Strategy
- Trend analysis: 
   If price is above center line → Uptrend
   If price is below center line → Downtrend

- Entry Point: 
   Enter trade when price exceeds 3 standard deviations from the center line

- Stop Loss: 
   Set stop loss based on center line

- Take Profit:
   Set take profit at 2x distance from stop loss

3. Special Strategies
- Bollinger Band Squeeze: 
   Enter when bands are narrow (low volatility), and breakout occurs in either direction

- Support/Resistance:
   When band breaks, confirm support or resistance level breach before entering

- Stop Loss: 
   Set stop loss based on center line

- Take Profit: 
   Set take profit at 2x distance from stop loss

4. Candlestick Pattern-Based Strategy
- Candlestick 1: Touching the upper band
- Candlestick 2: Touching the lower band
- Candlestick 3: Touching the center line and then rising/falling

- Enter when the trend line from Candlestick 1 and 3 is broken

5. Risk Management
- Always set a stop loss line and follow it

- Be cautious of increased risk when using leverage

6. Take Profit Tips
- After reaching 1st target, sell a portion and adjust stop loss to break even
- If price closes below center line, sell remaining position
`

// NewBitcoinStrategyPrompt ...
func NewBitcoinStrategyPrompt() string {
	return BitcoinStrategyTemplate
}
