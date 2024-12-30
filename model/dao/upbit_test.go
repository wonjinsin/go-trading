package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUpbitTokenPayload(t *testing.T) {
	// Given
	accessKey := "test-access-key"

	// When
	payload := NewUpbitTokenPayload(accessKey)

	// Then
	assert.Equal(t, accessKey, payload.AccessKey)
	assert.NotEmpty(t, payload.Nonce)
}

func TestNewSHA512UpbitTokenPayload(t *testing.T) {
	// Given
	accessKey := "test-access-key"
	query := "market=KRW-BTC&side=bid&price=1000&ord_type=price"

	// When
	payload := NewSHA512UpbitTokenPayload(accessKey, query)

	// Then
	assert.Equal(t, accessKey, payload.AccessKey)
	assert.NotEmpty(t, payload.Nonce)
	assert.Equal(t, "SHA512", payload.QueryHashAlg)
	assert.NotEmpty(t, payload.QueryHash)
}

func TestUpbitTokenPayload_GenerateJWT(t *testing.T) {
	// Given
	payload := &UpbitTokenPayload{
		AccessKey: "test-access-key",
		Nonce:     "test-nonce",
	}
	secretKey := "test-secret-key"

	// When
	token, err := payload.GenerateJWT(secretKey)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestUpbitAccounts_GetAccountByCurrency(t *testing.T) {
	// Given
	accounts := UpbitAccounts{
		{Currency: UpbitCurrencyKRW, Balance: "1000000"},
		{Currency: UpbitCurrencyBTC, Balance: "1.5"},
	}

	// When
	krwAccount := accounts.GetAccountByCurrency(UpbitCurrencyKRW)
	btcAccount := accounts.GetAccountByCurrency(UpbitCurrencyBTC)
	nonExistentAccount := accounts.GetAccountByCurrency("NON_EXISTENT")

	// Then
	assert.NotNil(t, krwAccount)
	assert.Equal(t, "1000000", krwAccount.Balance)

	assert.NotNil(t, btcAccount)
	assert.Equal(t, "1.5", btcAccount.Balance)

	assert.Nil(t, nonExistentAccount)
}

func TestNewUpbitOrderBuy(t *testing.T) {
	// Given
	market := UpbitStockBTC
	price := uint64(50000000)

	// When
	order := NewUpbitOrderBuy(market, price)

	// Then
	assert.Equal(t, market, order.Market)
	assert.Equal(t, UpbitOrderSideBuy, order.Side)
	assert.Equal(t, "50000000", order.Price)
	assert.Equal(t, UpbitOrderTypePrice, order.OrderType)
	assert.NotEmpty(t, order.Identifier)
}

func TestUpbitOrderBuy_GetQuery(t *testing.T) {
	// Given
	order := &UpbitOrderBuy{
		Market:     UpbitStockBTC,
		Side:       UpbitOrderSideBuy,
		Price:      "50000000",
		OrderType:  UpbitOrderTypePrice,
		Identifier: "test-identifier",
	}

	// When
	query := order.GetQuery()

	// Then
	expected := "market=KRW-BTC&side=bid&price=50000000&ord_type=price&identifier=test-identifier"
	assert.Equal(t, expected, query)
}

func TestNewUpbitOrderSell(t *testing.T) {
	// Given
	market := UpbitStockBTC
	volume := 1.5

	// When
	order := NewUpbitOrderSell(market, volume)

	// Then
	assert.Equal(t, market, order.Market)
	assert.Equal(t, UpbitOrderSideSell, order.Side)
	assert.Equal(t, "1.5", order.Volume)
	assert.Equal(t, UpbitOrderTypeMarket, order.OrderType)
	assert.NotEmpty(t, order.Identifier)
}

func TestUpbitOrderSell_GetQuery(t *testing.T) {
	// Given
	order := &UpbitOrderSell{
		Market:     UpbitStockBTC,
		Side:       UpbitOrderSideSell,
		Volume:     "1.5",
		OrderType:  UpbitOrderTypeMarket,
		Identifier: "test-identifier",
	}

	// When
	query := order.GetQuery()

	// Then
	expected := "market=KRW-BTC&side=ask&volume=1.5&ord_type=market&identifier=test-identifier"
	assert.Equal(t, expected, query)
}
