package service

import (
	"context"
	"magmar/config"
	"magmar/repository"
	"magmar/util"
)

type qaUsecase struct {
	conf          *config.ViperConfig
	openAIqaRepo  repository.QaRepository
	upbitBankRepo repository.BankRepository
}

// NewQaService ...
func NewQaService(conf *config.ViperConfig, qaRepo repository.QaRepository, upbitBankRepo repository.BankRepository) QaService {
	return &qaUsecase{
		conf:          conf,
		openAIqaRepo:  qaRepo,
		upbitBankRepo: upbitBankRepo,
	}
}

// QaService ...
func (q *qaUsecase) Ask(ctx context.Context) (err error) {
	zlog.With(ctx).Infow(util.LogSvc)
	err = q.upbitBankRepo.Buy(ctx)
	if err != nil {
		zlog.With(ctx).Errorw("Get token failed", "err", err)
		return err
	}
	return nil
	// return q.openAIqaRepo.Ask(ctx)
}
