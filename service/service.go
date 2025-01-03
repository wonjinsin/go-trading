package service

import (
	"context"
	"log"
	"magmar/config"
	"magmar/model"
	"magmar/model/dto"
	"magmar/repository"
	"magmar/util"
	"os"
)

var zlog *util.Logger

func init() {
	var err error
	zlog, err = util.NewLogger()
	if err != nil {
		log.Fatalf("InitLog module[service] err[%s]", err.Error())
		os.Exit(1)
	}
}

// Init ...
func Init(conf *config.ViperConfig, repo *repository.Repository) (*Service, error) {
	dealSvc := NewDealService(conf,
		repo.OpenAIQa,
		repo.UpbitBank,
		repo.AlternativeGreed,
		repo.News,
		repo.Transaction,
	)
	bankSvc := NewBankService(
		repo.UpbitBank,
		repo.Transaction,
	)
	return &Service{
		Deal: dealSvc,
		Bank: bankSvc,
	}, nil
}

// Service ...
type Service struct {
	Deal DealService
	Bank BankService
}

// DealService ...
type DealService interface {
	Deal(ctx context.Context) (err error)
}

// BankService ...
type BankService interface {
	Deposit(ctx context.Context, deposit dto.Deposit) (*model.TransactionAggregate, error)
	Withdrawal(ctx context.Context, withdrawal dto.Withdrawal) (*model.TransactionAggregate, error)
}
