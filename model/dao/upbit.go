package dao

import (
	"crypto/sha512"
	"encoding/hex"
	"net/url"
	"strconv"

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
}

// NewUpbitTokenPayload ...
func NewUpbitTokenPayload(accessKey string) *UpbitTokenPayload {
	return &UpbitTokenPayload{
		AccessKey: accessKey,
		Nonce:     uuid.New().String(),
	}
}

// NewSHA512UpbitTokenPayload ...
func NewSHA512UpbitTokenPayload(accessKey string, hash string) *UpbitTokenPayload {
	return &UpbitTokenPayload{
		AccessKey:    accessKey,
		Nonce:        uuid.New().String(),
		QueryHash:    hash,
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
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim).
		SignedString(secretKeyByte)
}

// UpbitAccount ...
type UpbitAccount struct {
	Currency string `json:"currency"`
	Balance  string `json:"balance"`
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
	OrderType  UpbitOrderType `json:"ord_type"`
	Price      string         `json:"price"`
	Identifier string         `json:"identifier"`
}

// NewUpbitOrderBuy ...
func NewUpbitOrderBuy(market UpbitStock, price uint64) *UpbitOrderBuy {
	return &UpbitOrderBuy{
		Market:     market,
		Side:       UpbitOrderSideBuy,
		OrderType:  UpbitOrderTypePrice,
		Price:      strconv.FormatUint(price, 10),
		Identifier: uuid.New().String(),
	}
}

// ToSHA512 ...
func (p *UpbitOrderBuy) ToSHA512() string {
	query := url.Values{
		"market":     []string{string(p.Market)},
		"side":       []string{string(p.Side)},
		"price":      []string{p.Price},
		"ord_type":   []string{string(p.OrderType)},
		"identifier": []string{p.Identifier},
	}

	hash := sha512.New()
	hash.Write([]byte(query.Encode()))
	return hex.EncodeToString(hash.Sum(nil))
}
