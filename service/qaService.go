package service

import (
	"context"
	"magmar/config"
	"magmar/repository"
	"magmar/util"
)

type qaUsecase struct {
	conf         *config.ViperConfig
	openAIqaRepo repository.QaRepository
}

// NewQaService ...
func NewQaService(conf *config.ViperConfig, qaRepo repository.QaRepository) QaService {
	return &qaUsecase{
		conf:         conf,
		openAIqaRepo: qaRepo,
	}
}

// QaService ...
func (q *qaUsecase) Ask(ctx context.Context) (err error) {
	zlog.With(ctx).Infow(util.LogSvc)
	return q.openAIqaRepo.Ask(ctx)
}
