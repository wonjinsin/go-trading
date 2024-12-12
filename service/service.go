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
	qaSvc := NewQaService(conf, repo.OpenAIQa, repo.UpbitBank)
	return &Service{
		Qa: qaSvc,
	}, nil
}

// Service ...
type Service struct {
	Qa QaService
}

// QaService ...
type QaService interface {
	Ask(ctx context.Context) (err error)
}
