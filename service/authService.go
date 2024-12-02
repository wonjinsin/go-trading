package service

import (
	"context"
	"magmar/config"
	"magmar/repository"
)

type qaUsecase struct {
	conf   *config.ViperConfig
	qaRepo repository.QaRepository
}

// NewQaService ...
func NewQaService(conf *config.ViperConfig, qaRepo repository.QaRepository) QaService {
	return &qaUsecase{
		conf:   conf,
		qaRepo: qaRepo,
	}
}

// QaService ...
func (q *qaUsecase) Ask(ctx context.Context) (err error) {
	zlog.With(ctx).Infow("[New Service Request]", "ask", "ask")
	return q.qaRepo.Ask(ctx)
}
