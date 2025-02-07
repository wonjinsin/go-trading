package util

// Coin ...
type Coin string

// CoinConst ...
const (
	Balance Coin = "Balance"
	BitCoin Coin = "Bitcoin"
)

// Bank ...
type Bank string

// BankConst ...
const (
	Upbit         Bank = "Upbit"
	UpbitCurrency Bank = "UpbitCurrency"
	News          Bank = "News"
)

// CoinMap ...
var CoinMap = map[Coin]map[Bank]string{
	BitCoin: {Upbit: "BTC-KRW", UpbitCurrency: "BTC", News: "bitcoin"},
	Balance: {Upbit: "KRW", UpbitCurrency: "KRW", News: "korea"},
}
