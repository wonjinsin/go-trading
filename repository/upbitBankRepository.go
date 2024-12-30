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
func NewUpbitBankRepository(conf *config.ViperConfig) BankRepository {
	return &upbitBankRepository{
		accessKey: conf.GetString(util.UpbitAccessKey),
		secretKey: conf.GetString(util.UpbitSecretKey),
		conn:      resty.New(),
		apiURL:    util.APIURLUpbit,
	}
}

// GetOrderBook ...
func (b *upbitBankRepository) GetOrderBook(ctx context.Context, stock dao.UpbitStock) (orderBook *model.OrderBook, err error) {
	zlog.With(ctx).Infow(util.LogRepo)
	resp, err := b.conn.R().
		SetResult(&orderBook).
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

	return orderBook, nil
}

// GetMarketPriceData ...
func (b *upbitBankRepository) GetMarketPriceData(ctx context.Context, stock dao.UpbitStock, date uint) (marketPrices model.MarketPrices, err error) {
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

	account := accounts.GetAccountByCurrency(dao.UpbitCurrencyKRW)
	if account == nil {
		zlog.With(ctx).Errorw("Get account currecy not exist", "accounts", accounts)
		return nil, errors.NotFoundf("Get account currecy not exist")
	}

	balance := &model.BankBalance{
		Currency: string(account.Currency),
		Balance:  util.ParseFloat64(account.Balance),
	}

	return balance, nil
}

// GetBitCoinBalance ...
func (b *upbitBankRepository) GetBitCoinBalance(ctx context.Context) (*model.BankBalance, error) {
	zlog.With(ctx).Infow(util.LogRepo)
	accounts, err := b.getBalance(ctx)
	if err != nil {
		zlog.With(ctx).Errorw("Get accounts failed", "err", err)
		return nil, err
	}

	account := accounts.GetAccountByCurrency(dao.UpbitCurrencyBTC)
	if account == nil {
		zlog.With(ctx).Errorw("Get account currecy not exist", "accounts", accounts)
		return nil, errors.NotFoundf("Get account currecy not exist")
	}

	balance := &model.BankBalance{
		Currency: string(account.Currency),
		Balance:  util.ParseFloat64(account.Balance),
	}

	return balance, nil
}

// Buy ...
func (b *upbitBankRepository) Buy(ctx context.Context, amount uint64) (err error) {
	zlog.With(ctx).Infow(util.LogRepo, "amount", amount)
	orderBuy := dao.NewUpbitOrderBuy(dao.UpbitStockBTC, amount)
	zlog.With(ctx).Infow("Order calculated", "orderBuy", orderBuy)

	token, err := b.getSHA512Token(orderBuy)
	if err != nil {
		zlog.With(ctx).Errorw("Generate JWT failed", "err", err)
		return err
	}

	resp, err := b.conn.R().
		SetAuthToken(token).
		SetBody(orderBuy).
		Post(fmt.Sprintf("%s/v1/orders", b.apiURL))
	if err != nil {
		zlog.With(ctx).Errorw("Buy failed", "status", resp.StatusCode(), "resp", resp.String(), "err", err)
		return err
	}

	if resp.StatusCode() != http.StatusCreated {
		zlog.With(ctx).Errorw("Buy failed", "status", resp.StatusCode(), "resp", resp.String())
		return errors.NotImplementedf("Buy failed")
	}

	return nil
}

// Sell ...
func (b *upbitBankRepository) Sell(ctx context.Context, amount float64) (err error) {
	zlog.With(ctx).Infow(util.LogRepo, "amount", amount)
	orderSell := dao.NewUpbitOrderSell(dao.UpbitStockBTC, amount)
	zlog.With(ctx).Infow("Order calculated", "orderSell", orderSell)

	token, err := b.getSHA512Token(orderSell)
	if err != nil {
		zlog.With(ctx).Errorw("Generate JWT failed", "err", err)
		return err
	}

	resp, err := b.conn.R().
		SetAuthToken(token).
		SetBody(orderSell).
		Post(fmt.Sprintf("%s/v1/orders", b.apiURL))
	if err != nil {
		zlog.With(ctx).Errorw("Sell failed", "status", resp.StatusCode(), "resp", resp.String(), "err", err)
		return err
	}

	if resp.StatusCode() != http.StatusCreated {
		zlog.With(ctx).Errorw("Buy failed", "status", resp.StatusCode(), "resp", resp.String())
		return errors.NotImplementedf("Buy failed")
	}

	return nil
}

func (b *upbitBankRepository) getToken() (token string, err error) {
	tokenPayload := dao.NewUpbitTokenPayload(b.accessKey)
	return tokenPayload.GenerateJWT(b.secretKey)
}

func (b *upbitBankRepository) getSHA512Token(queryable dao.Queryable) (token string, err error) {
	tokenPayload := dao.NewSHA512UpbitTokenPayload(b.accessKey, queryable.GetQuery())
	return tokenPayload.GenerateJWT(b.secretKey)
}
