package service

import (
	"context"
	"magmar/config"
	"magmar/model"
	"magmar/model/dao"
	"magmar/prompt"
	"magmar/repository"
	"magmar/util"
)

type dealUsecase struct {
	conf          *config.ViperConfig
	feePercent    uint
	feeScale      uint
	openAIqaRepo  repository.QaRepository
	upbitBankRepo repository.BankRepository
}

// NewDealService ...
func NewDealService(conf *config.ViperConfig, qaRepo repository.QaRepository, upbitBankRepo repository.BankRepository) DealService {
	return &dealUsecase{
		conf:          conf,
		feePercent:    conf.GetUint(util.UpbitFeePercent),
		feeScale:      conf.GetUint(util.UpbitFeeScale),
		openAIqaRepo:  qaRepo,
		upbitBankRepo: upbitBankRepo,
	}
}

// Deal ...
func (d *dealUsecase) Deal(ctx context.Context) (err error) {
	zlog.With(ctx).Infow(util.LogSvc)
	decision, err := d.ask(ctx)
	if err != nil {
		zlog.With(ctx).Errorw("Ask failed", "err", err)
		return err
	}

	zlog.With(ctx).Infow("Got decision", "decision", decision)

	switch decision.Decision {
	case model.DecisionStateBuy:
		err = d.buy(ctx)
	case model.DecisionStateSell:
		err = d.sell(ctx)
	case model.DecisionStateHold:
		err = nil
	}

	if err != nil {
		zlog.With(ctx).Warnw("Buy failed", "err", err)
		return err
	}

	zlog.With(ctx).Infow("Process done")
	return nil
}

func (d *dealUsecase) ask(ctx context.Context) (decision *model.Decision, err error) {
	marketPrices, err := d.upbitBankRepo.GetMarketPriceData(ctx, dao.UpbitStockBTC, 60)
	if err != nil {
		zlog.With(ctx).Warnw("Get token failed", "err", err)
		return nil, err
	}

	decision, err = d.openAIqaRepo.Ask(ctx, prompt.NewBitcoinPrompt(marketPrices))
	if err != nil {
		zlog.With(ctx).Warnw("Ask failed", "err", err)
		return nil, err
	}

	return decision, nil
}

func (d *dealUsecase) buy(ctx context.Context) (err error) {
	zlog.With(ctx).Infow("Buy start")

	balance, err := d.upbitBankRepo.GetBalance(ctx)
	if err != nil {
		zlog.With(ctx).Warnw("Get balance failed", "err", err)
		return err
	}

	zlog.With(ctx).Infow("Got balance", "balance", balance)
	amount := balance.GetBuyAmount(d.feePercent, d.feeScale)

	err = d.upbitBankRepo.Buy(ctx, amount)
	if err != nil {
		zlog.With(ctx).Warnw("Buy failed", "err", err)
		return err
	}

	return nil
}

func (d *dealUsecase) sell(_ context.Context) (err error) {
	return nil
}
