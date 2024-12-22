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

// GetBalance ...
func (b *upbitBankRepository) GetBalance(ctx context.Context) (*model.BankBalance, error) {
	zlog.With(ctx).Infow(util.LogRepo)
	token, err := b.getToken()
	if err != nil {
		zlog.With(ctx).Errorw("Generate JWT failed", "err", err)
		return nil, err
	}

	var accounts []*dao.UpbitAccount
	resp, err := b.conn.R().
		SetHeader("Authorization", token).
		SetResult(&accounts).
		Get(fmt.Sprintf("%s/v1/accounts", b.apiURL))
	if err != nil {
		zlog.With(ctx).Errorw("Get accounts failed", "err", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		zlog.With(ctx).Errorw("Get accounts failed", "status", resp.StatusCode())
		return nil, errors.NotImplementedf("Get accounts failed")
	}
	account := accounts[0]
	balance := &model.BankBalance{
		Currency: account.Currency,
		Balance:  util.ParseUint64(account.Balance),
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

func (b *upbitBankRepository) getToken() (token string, err error) {
	tokenPayload := dao.NewUpbitTokenPayload(b.accessKey)
	return tokenPayload.GenerateJWT(b.secretKey)
}

func (b *upbitBankRepository) getSHA512Token(queryable dao.Queryable) (token string, err error) {
	tokenPayload := dao.NewSHA512UpbitTokenPayload(b.accessKey, queryable.GetQuery())
	return tokenPayload.GenerateJWT(b.secretKey)
}
