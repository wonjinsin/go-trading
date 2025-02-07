package repository

import (
	"context"
	"fmt"
	"magmar/config"
	"magmar/model"
	"magmar/model/dao"
	"magmar/util"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/juju/errors"
)

type upbitBankRepository struct {
	accessKey string
	secretKey string
	conn      *resty.Client
	apiURL    util.APIURL
}

// NewUpbitBankRepository ...
func NewUpbitBankRepository(conf *config.ViperConfig) StockBankRepository {
	return &upbitBankRepository{
		accessKey: conf.GetString(util.UpbitAccessKey),
		secretKey: conf.GetString(util.UpbitSecretKey),
		conn:      resty.New(),
		apiURL:    util.APIURLUpbit,
	}
}

// GetOrderBook ...
func (b *upbitBankRepository) GetOrderBooks(ctx context.Context, stock model.StockName) (orderBooks model.OrderBooks, err error) {
	zlog.With(ctx).Infow(util.LogRepo)
	resp, err := b.conn.R().
		SetResult(&orderBooks).
		SetQueryParam("level", "0").
		SetQueryParam("markets", string(stock)).
		Get(fmt.Sprintf("%s/v1/orderbook", b.apiURL))
	if err != nil {
		zlog.With(ctx).Errorw("Get order book failed", "err", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		zlog.With(ctx).Errorw("Get order book failed", "status", resp.StatusCode())
		return nil, errors.NotImplementedf("Get order book failed")
	}

	return orderBooks, nil
}

// GetMarketPriceDataDay ...
func (b *upbitBankRepository) GetMarketPriceDataDay(ctx context.Context, stock model.StockName, date uint) (marketPrices model.MarketPrices, err error) {
	zlog.With(ctx).Infow(util.LogRepo)
	var upbitMarketPrices dao.UpbitMarketPrices
	resp, err := b.conn.R().
		SetResult(&upbitMarketPrices).
		SetQueryParam("market", string(stock)).
		SetQueryParam("count", fmt.Sprintf("%d", date)).
		Get(fmt.Sprintf("%s/v1/candles/days", b.apiURL))
	if err != nil {
		zlog.With(ctx).Errorw("Get market price failed", "err", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		zlog.With(ctx).Errorw("Get market price failed", "status", resp.StatusCode())
		return nil, errors.NotImplementedf("Get market price failed")
	}

	marketPrices = model.NewMarketPriceByUpbit(upbitMarketPrices)
	marketPrices.SetIndicators()
	return marketPrices, nil
}

// GetMarketPriceDataMin ...
func (b *upbitBankRepository) GetMarketPriceDataMin(ctx context.Context, stock model.StockName, interval uint) (marketPrices model.MarketPrices, err error) {
	zlog.With(ctx).Infow(util.LogRepo)
	var upbitMarketPrices dao.UpbitMarketPrices
	resp, err := b.conn.R().
		SetResult(&upbitMarketPrices).
		SetQueryParam("market", string(stock)).
		SetQueryParam("count", fmt.Sprintf("%d", 1440/interval)).
		Get(fmt.Sprintf("%s/v1/candles/minutes/%d", b.apiURL, interval))
	if err != nil {
		zlog.With(ctx).Errorw("Get market price failed", "err", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		zlog.With(ctx).Errorw("Get market price failed", "status", resp.StatusCode())
		return nil, errors.NotImplementedf("Get market price failed")
	}

	marketPrices = model.NewMarketPriceByUpbit(upbitMarketPrices)
	marketPrices.SetIndicators()
	return marketPrices, nil
}

func (b *upbitBankRepository) getBalance(ctx context.Context) (dao.UpbitAccounts, error) {
	token, err := b.getToken()
	if err != nil {
		zlog.With(ctx).Errorw("Generate JWT failed", "err", err)
		return nil, err
	}

	var accounts dao.UpbitAccounts
	resp, err := b.conn.R().
		SetHeader("Authorization", token).
		SetResult(&accounts).
		Get(fmt.Sprintf("%s/v1/accounts", b.apiURL))
	if err != nil {
		zlog.With(ctx).Errorw("Get accounts failed", "err", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		zlog.With(ctx).Errorw("Get accounts failed", "status", resp.StatusCode(), "resp", resp.String())
		return nil, errors.NotImplementedf("Get accounts failed")
	}

	return accounts, nil
}

// GetBalance ...
func (b *upbitBankRepository) GetBalance(ctx context.Context) (*model.BankBalance, error) {
	zlog.With(ctx).Infow(util.LogRepo)
	accounts, err := b.getBalance(ctx)
	if err != nil {
		zlog.With(ctx).Errorw("Get accounts failed", "err", err)
		return nil, err
	}

	account := accounts.GetAccountByCurrency(util.CoinMap[util.Balance][util.UpbitCurrency])
	if account == nil {
		zlog.With(ctx).Errorw("Get account currecy not exist", "accounts", accounts)
		return nil, errors.NotFoundf("Get account currecy not exist")
	}

	balance := &model.BankBalance{
		Currency:    string(account.Currency),
		Balance:     util.ParseFloat64(account.Balance),
		AvgBuyPrice: util.ParseFloat64("1"),
	}

	return balance, nil
}

// GetCoinBalance ...
func (b *upbitBankRepository) GetCoinBalance(ctx context.Context, currency string) (*model.BankBalance, error) {
	zlog.With(ctx).Infow(util.LogRepo)
	accounts, err := b.getBalance(ctx)
	if err != nil {
		zlog.With(ctx).Errorw("Get accounts failed", "err", err)
		return nil, err
	}

	account := accounts.GetAccountByCurrency(currency)
	if account == nil {
		zlog.With(ctx).Errorw("Get account currecy not exist", "accounts", accounts)
		return nil, errors.NotFoundf("Get account currecy not exist")
	}

	balance := &model.BankBalance{
		Currency:    string(account.Currency),
		Balance:     util.ParseFloat64(account.Balance),
		AvgBuyPrice: util.ParseFloat64(account.AvgBuyPrice),
	}

	return balance, nil
}

// Buy ...
func (b *upbitBankRepository) Buy(ctx context.Context, amount uint64) (result *model.BankTransactionResult, err error) {
	zlog.With(ctx).Infow(util.LogRepo, "amount", amount)
	orderBuy := dao.NewUpbitOrderBuy(dao.UpbitStockBTC, amount)
	zlog.With(ctx).Infow("Order calculated", "orderBuy", orderBuy)

	token, err := b.getSHA512Token(orderBuy)
	if err != nil {
		zlog.With(ctx).Errorw("Generate JWT failed", "err", err)
		return nil, err
	}

	var trResult *dao.UpbitTransactionResult
	resp, err := b.conn.R().
		SetAuthToken(token).
		SetResult(&trResult).
		SetBody(orderBuy).
		Post(fmt.Sprintf("%s/v1/orders", b.apiURL))
	if err != nil {
		zlog.With(ctx).Errorw("Buy failed", "status", resp.StatusCode(), "resp", resp.String(), "err", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusCreated {
		zlog.With(ctx).Errorw("Buy failed", "status", resp.StatusCode(), "resp", resp.String())
		return nil, errors.NotImplementedf("Buy failed")
	}

	result = model.NewBankTransactionResultBuy(trResult)
	return result, nil
}

// Sell ...
func (b *upbitBankRepository) Sell(ctx context.Context, amount float64) (result *model.BankTransactionResult, err error) {
	zlog.With(ctx).Infow(util.LogRepo, "amount", amount)
	orderSell := dao.NewUpbitOrderSell(dao.UpbitStockBTC, amount)
	zlog.With(ctx).Infow("Order calculated", "orderSell", orderSell)

	token, err := b.getSHA512Token(orderSell)
	if err != nil {
		zlog.With(ctx).Errorw("Generate JWT failed", "err", err)
		return nil, err
	}

	var trResult *dao.UpbitTransactionResult
	resp, err := b.conn.R().
		SetAuthToken(token).
		SetResult(&trResult).
		SetBody(orderSell).
		Post(fmt.Sprintf("%s/v1/orders", b.apiURL))
	if err != nil {
		zlog.With(ctx).Errorw("Sell failed", "status", resp.StatusCode(), "resp", resp.String(), "err", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusCreated {
		zlog.With(ctx).Errorw("Buy failed", "status", resp.StatusCode(), "resp", resp.String())
		return nil, errors.NotImplementedf("Buy failed")
	}

	result = model.NewBankTransactionResultSell(trResult)
	return result, nil
}

func (b *upbitBankRepository) getToken() (token string, err error) {
	tokenPayload := dao.NewUpbitTokenPayload(b.accessKey)
	return tokenPayload.GenerateJWT(b.secretKey)
}

func (b *upbitBankRepository) getSHA512Token(queryable dao.Queryable) (token string, err error) {
	tokenPayload := dao.NewSHA512UpbitTokenPayload(b.accessKey, queryable.GetQuery())
	return tokenPayload.GenerateJWT(b.secretKey)
}
