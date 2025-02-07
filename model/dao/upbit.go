package dao

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UpbitStock ...
type UpbitStock string

// UpbitStockConst ...
const (
	UpbitStockBTC UpbitStock = "KRW-BTC"
)

// UpbitOrderSide ...
type UpbitOrderSide string

// UpbitOrderSideConst ...
const (
	UpbitOrderSideBuy  UpbitOrderSide = "bid"
	UpbitOrderSideSell UpbitOrderSide = "ask"
)

// UpbitOrderType ...
type UpbitOrderType string

// UpbitOrderTypeConst ...
const (
	UpbitOrderTypePrice  UpbitOrderType = "price"
	UpbitOrderTypeMarket UpbitOrderType = "market"
)

// UpbitTokenPayload ...
type UpbitTokenPayload struct {
	AccessKey    string `json:"access_key"`
	Nonce        string `json:"nonce"`
	QueryHash    string `json:"query_hash"`
	QueryHashAlg string `json:"query_hash_alg"`
	Identifier   string `json:"identifier"`
}

// NewUpbitTokenPayload ...
func NewUpbitTokenPayload(accessKey string) *UpbitTokenPayload {
	return &UpbitTokenPayload{
		AccessKey: accessKey,
		Nonce:     uuid.New().String(),
	}
}

// NewSHA512UpbitTokenPayload ...
func NewSHA512UpbitTokenPayload(accessKey string, query string) *UpbitTokenPayload {
	hash := sha512.New()
	hash.Write([]byte(query))
	hashString := hex.EncodeToString(hash.Sum(nil))

	return &UpbitTokenPayload{
		AccessKey:    accessKey,
		Nonce:        uuid.New().String(),
		QueryHash:    hashString,
		QueryHashAlg: "SHA512",
	}
}

// GenerateJWT ...
func (p *UpbitTokenPayload) GenerateJWT(secretKey string) (string, error) {
	secretKeyByte := []byte(secretKey)
	tokenClaim := jwt.MapClaims{
		"access_key":     p.AccessKey,
		"nonce":          p.Nonce,
		"query_hash":     p.QueryHash,
		"query_hash_alg": p.QueryHashAlg,
		"identifier":     p.Identifier,
		"iat":            time.Now().Add(time.Hour * 1).Unix(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim).
		SignedString(secretKeyByte)
}

// UpbitAccount ...
type UpbitAccount struct {
	Currency    string `json:"currency"`
	Balance     string `json:"balance"`
	AvgBuyPrice string `json:"avg_buy_price"`
}

// UpbitAccounts ...
type UpbitAccounts []*UpbitAccount

// GetAccountByCurrency ...
func (us UpbitAccounts) GetAccountByCurrency(currency string) *UpbitAccount {
	for _, account := range us {
		if account.Currency == currency {
			return account
		}
	}
	return nil
}

// UpbitMarketPrice ...
type UpbitMarketPrice struct {
	Market            UpbitStock `json:"market"`
	CandleDateTimeUTC string     `json:"candle_date_time_utc"`
	TradePrice        float64    `json:"trade_price"`
}

// UpbitMarketPrices ...
type UpbitMarketPrices []*UpbitMarketPrice

// UpbitOrderBuy ...
type UpbitOrderBuy struct {
	Market     UpbitStock     `json:"market"`
	Side       UpbitOrderSide `json:"side"`
	Price      string         `json:"price"`
	OrderType  UpbitOrderType `json:"ord_type"`
	Identifier string         `json:"identifier"`
}

// NewUpbitOrderBuy ...
func NewUpbitOrderBuy(market UpbitStock, price uint64) *UpbitOrderBuy {
	return &UpbitOrderBuy{
		Market:     market,
		Side:       UpbitOrderSideBuy,
		Price:      strconv.FormatUint(price, 10),
		OrderType:  UpbitOrderTypePrice,
		Identifier: uuid.New().String(),
	}
}

// GetQuery ...
func (p *UpbitOrderBuy) GetQuery() string {
	return fmt.Sprintf("market=%s&side=%s&price=%s&ord_type=%s&identifier=%s",
		p.Market,
		p.Side,
		p.Price,
		p.OrderType,
		p.Identifier,
	)
}

// UpbitOrderSell ...
type UpbitOrderSell struct {
	Market     UpbitStock     `json:"market"`
	Side       UpbitOrderSide `json:"side"`
	Volume     string         `json:"volume"`
	OrderType  UpbitOrderType `json:"ord_type"`
	Identifier string         `json:"identifier"`
}

// NewUpbitOrderSell ...
func NewUpbitOrderSell(market UpbitStock, price float64) *UpbitOrderSell {
	return &UpbitOrderSell{
		Market:     market,
		Side:       UpbitOrderSideSell,
		Volume:     strconv.FormatFloat(price, 'f', -1, 64),
		OrderType:  UpbitOrderTypeMarket,
		Identifier: uuid.New().String(),
	}
}

// GetQuery ...
func (p *UpbitOrderSell) GetQuery() string {
	return fmt.Sprintf("market=%s&side=%s&volume=%s&ord_type=%s&identifier=%s",
		p.Market,
		p.Side,
		p.Volume,
		p.OrderType,
		p.Identifier,
	)
}

// UpbitTransactionResult ...
type UpbitTransactionResult struct {
	Market     UpbitStock `json:"market"`
	Identifier *string    `json:"identifier"`
	Price      *string    `json:"price"`
	Volume     *string    `json:"volume"`
}
