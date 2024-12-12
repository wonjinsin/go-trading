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
)

type upbitBankRepository struct {
	accessKey string
	secretKey string
	conn      *resty.Client
}

// NewUpbitBankRepository ...
func NewUpbitBankRepository(conf *config.ViperConfig) BankRepository {
	return &upbitBankRepository{
		accessKey: conf.GetString(util.UpbitAccessKey),
		secretKey: conf.GetString(util.UpbitSecretKey),
		conn:      resty.New(),
	}
}

// Buy ...
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
		Get("https://api.upbit.com/v1/accounts")
	if err != nil {
		zlog.With(ctx).Errorw("Get accounts failed", "err", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		zlog.With(ctx).Errorw("Get accounts failed", "status", resp.StatusCode())
		return nil, err
	}
	account := accounts[0]
	balance := &model.BankBalance{
		Currency: account.Currency,
		Balance:  account.Balance,
	}

	return balance, nil
}

// Buy ...
func (b *upbitBankRepository) Buy(ctx context.Context) (err error) {
	zlog.With(ctx).Infow(util.LogRepo)
	token, err := b.getToken()
	if err != nil {
		zlog.With(ctx).Errorw("Generate JWT failed", "err", err)
		return err
	}

	resp, err := b.conn.R().
		SetHeader("Authorization", token).
		Get("https://api.upbit.com/v1/accounts")
	if err != nil {
		zlog.With(ctx).Errorw("Get accounts failed", "err", err)
		return err
	}

	fmt.Println(resp.String())
	return nil
}

func (b *upbitBankRepository) getToken() (token string, err error) {
	tokenPayload := dao.NewUpbitTokenPayload(b.accessKey)
	return tokenPayload.GenerateJWT(b.secretKey)
}
