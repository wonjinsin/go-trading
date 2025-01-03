package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"magmar/config"
	"magmar/model"
	"magmar/util"

	"github.com/tmc/langchaingo/llms/openai"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	mysqlConn, err := mysqlConnect(magmar)
	if err != nil {
		return nil, err
	}

	openAPIConn, err := openAPIConnect(magmar)
	if err != nil {
		return nil, err
	}

	db := &model.DB{
		OpenAI: openAPIConn,
		MainDB: mysqlConn,
	}

	qaRepo := NewOpenAPIQaRepository(db.OpenAI)
	upbitBankRepo := NewUpbitBankRepository(magmar)
	alternativeGreedRepo := NewAlternativeGreedRepository()
	newsAPIRepo := NewNewsAPIRepository(magmar)
	transactionRepo := NewGormTransactionRepository(db.MainDB)

	return &Repository{
		OpenAIQa:         qaRepo,
		UpbitBank:        upbitBankRepo,
		AlternativeGreed: alternativeGreedRepo,
		News:             newsAPIRepo,
		Transaction:      transactionRepo,
	}, nil
}

// Repository ...
type Repository struct {
	OpenAIQa         QaRepository
	UpbitBank        StockBankRepository
	AlternativeGreed GreedRepository
	News             NewsRepository
	Transaction      TransactionRepository
}

func openAPIConnect(magmar *config.ViperConfig) (*openai.LLM, error) {
	opt := openai.WithToken(magmar.GetString(util.OpenAPIKey))
	return openai.New(opt)
}

func mysqlConnect(magmar *config.ViperConfig) (*gorm.DB, error) {
	return gorm.Open(getDialector(magmar), getConfig())
}

func getDialector(magmar *config.ViperConfig) gorm.Dialector {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=UTC",
		magmar.GetString(util.DBUserKey),
		magmar.GetString(util.DBPasswordKey),
		magmar.GetString(util.DBHostKey),
		magmar.GetInt(util.DBPortKey),
		magmar.GetString(util.DBNameKey),
	)

	return mysql.Open(dbURI)
}

func getConfig() (gConfig *gorm.Config) {
	dbLogger := &dbLogger{zlog}
	gConfig = &gorm.Config{
		Logger:                                   dbLogger,
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	return gConfig
}

// QaRepository ...
type QaRepository interface {
	Ask(ctx context.Context, prompt string) (decision *model.Decision, err error)
}

// StockBankRepository ...
type StockBankRepository interface {
	GetOrderBooks(ctx context.Context, stock model.StockName) (orderBooks model.OrderBooks, err error)
	GetMarketPriceDataDay(ctx context.Context, stock model.StockName, date uint) (marketPrices model.MarketPrices, err error)
	GetMarketPriceDataMin(ctx context.Context, stock model.StockName, interval uint) (marketPrices model.MarketPrices, err error)
	GetBalance(ctx context.Context) (*model.BankBalance, error)
	GetBitCoinBalance(ctx context.Context) (*model.BankBalance, error)
	Buy(ctx context.Context, amount uint64) (*model.BankTransactionResult, error)
	Sell(ctx context.Context, amount float64) (*model.BankTransactionResult, error)
}

// GreedRepository ...
type GreedRepository interface {
	GetGreedIndex(ctx context.Context) (index *model.GreedIndex, err error)
}

// NewsRepository ...
type NewsRepository interface {
	GetNews(ctx context.Context, keywords []string) (newses model.Newses, err error)
}

// TransactionRepository ...
type TransactionRepository interface {
	NewTransaction(ctx context.Context, transaction *model.TransactionAggregate) (*model.TransactionAggregate, error)
	GetTransactions(ctx context.Context) (model.TransactionAggregates, error)
	GetTotalDeposit(ctx context.Context) (float64, error)
	GetTotalWithdrawal(ctx context.Context) (float64, error)
}
