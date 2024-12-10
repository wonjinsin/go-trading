package repository

import (
	"context"
	"magmar/config"
	"magmar/util"

	"github.com/go-resty/resty/v2"
)

type bankRepository struct {
	accessKey string
	secretKey string
	conn      *resty.Client
}

// NewUpbitBankRepository ...
func NewUpbitBankRepository(conf *config.ViperConfig) BankRepository {
	return &bankRepository{
		accessKey: conf.GetString(util.UpbitAccessKey),
		secretKey: conf.GetString(util.UpbitSecretKey),
		conn:      resty.New(),
	}
}

// Buy ...
func (b *bankRepository) Buy(ctx context.Context) (err error) {
	zlog.With(ctx).Infow(util.LogRepo)
	return nil
}
