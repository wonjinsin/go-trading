package service

import (
	"context"
	"log"
	"magmar/config"
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
	)
	return &Service{
		Deal: dealSvc,
	}, nil
}

// Service ...
type Service struct {
	Deal DealService
}

// DealService ...
type DealService interface {
	Deal(ctx context.Context) (err error)
}
