package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"magmar/config"
	"magmar/model"
	"magmar/model/dao"
	"magmar/util"

	"github.com/tmc/langchaingo/llms/openai"
	"gorm.io/gorm/logger"
)

var zlog *util.Logger

type dbLogger struct {
	*util.Logger
}

func (dl *dbLogger) LogMode(l logger.LogLevel) logger.Interface {
	return dl
}

func (dl *dbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	dl.Logger.With(ctx).Info(msg, data)
}

func (dl *dbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	dl.Logger.With(ctx).Warn(msg, data)
}

func (dl *dbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	dl.Logger.With(ctx).Error(msg, data)
}

func (dl *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		dl.Logger.With(ctx).Infow(err.Error(), "elapsed", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6), "rows", rows, "sql", sql)
	} else {
		dl.Logger.With(ctx).Infow("", "elapsed", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6), "rows", rows, "sql", sql)
	}
}

func init() {
	var err error
	zlog, err = util.NewLogger()
	if err != nil {
		log.Fatalf("InitLog module[service] err[%s]", err.Error())
		os.Exit(1)
	}
}

// Init ...
func Init(magmar *config.ViperConfig) (*Repository, error) {
	openAPIConn, err := openAPIConnect(magmar)
	if err != nil {
		return nil, err
	}

	db := &model.DB{
		OpenAI: openAPIConn,
	}

	qaRepo := NewOpenAPIQaRepository(db.OpenAI)
	upbitBankRepo := NewUpbitBankRepository(magmar)

	return &Repository{
		OpenAIQa:  qaRepo,
		UpbitBank: upbitBankRepo,
	}, nil
}

// Repository ...
type Repository struct {
	OpenAIQa  QaRepository
	UpbitBank BankRepository
}

func openAPIConnect(magmar *config.ViperConfig) (*openai.LLM, error) {
	opt := openai.WithToken(magmar.GetString(util.OpenAPIKey))
	return openai.New(opt)
}

// QaRepository ...
type QaRepository interface {
	Ask(ctx context.Context, prompt string) (decision *model.Decision, err error)
}

// BankRepository ...
type BankRepository interface {
	GetMarketPriceData(ctx context.Context, stock dao.UpbitStock, date uint) (marketPrices model.MarketPrices, err error)
	GetBalance(ctx context.Context) (*model.BankBalance, error)
	GetBitCoinBalance(ctx context.Context) (*model.BankBalance, error)
	Buy(ctx context.Context, amount uint64) (err error)
	Sell(ctx context.Context, amount float64) (err error)
}
