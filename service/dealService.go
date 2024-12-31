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
	minBuyAmount  uint64
	openAIqaRepo  repository.QaRepository
	upbitBankRepo repository.BankRepository
	greedRepo     repository.GreedRepository
	newsRepo      repository.NewsRepository
}

// NewDealService ...
func NewDealService(conf *config.ViperConfig,
	qaRepo repository.QaRepository,
	upbitBankRepo repository.BankRepository,
	greedRepo repository.GreedRepository,
	newsRepo repository.NewsRepository) DealService {
	return &dealUsecase{
		conf:          conf,
		feePercent:    conf.GetUint(util.UpbitFeePercent),
		feeScale:      conf.GetUint(util.UpbitFeeScale),
		minBuyAmount:  conf.GetUint64(util.UpbitMinBuyAmount),
		openAIqaRepo:  qaRepo,
		upbitBankRepo: upbitBankRepo,
		greedRepo:     greedRepo,
		newsRepo:      newsRepo,
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

	// todo: add decision percentage
	// todo: add to database result
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

// todo: add recursion asking also
func (d *dealUsecase) ask(ctx context.Context) (decision *model.Decision, err error) {
	marketPricesDay, err := d.upbitBankRepo.GetMarketPriceDataDay(ctx, dao.UpbitStockBTC, 60)
	if err != nil {
		zlog.With(ctx).Warnw("Get market price data day failed", "err", err)
		return nil, err
	}

	marketPricesMin, err := d.upbitBankRepo.GetMarketPriceDataMin(ctx, dao.UpbitStockBTC, 60)
	if err != nil {
		zlog.With(ctx).Warnw("Get market price data min failed", "err", err)
		return nil, err
	}

	orderBooks, err := d.upbitBankRepo.GetOrderBooks(ctx, dao.UpbitStockBTC)
	if err != nil {
		zlog.With(ctx).Warnw("Get order book failed", "err", err)
		return nil, err
	}

	// todo: add fear and greed index
	greedIndex, err := d.greedRepo.GetGreedIndex(ctx)
	if err != nil {
		zlog.With(ctx).Warnw("Get greed index failed", "err", err)
		return nil, err
	}

	news, err := d.newsRepo.GetNews(ctx, []string{"bitcoin", "predict"})
	if err != nil {
		zlog.With(ctx).Warnw("Get news failed", "err", err)
		return nil, err
	}

	decision, err = d.openAIqaRepo.Ask(ctx, prompt.NewBitcoinPrompt(
		marketPricesDay,
		marketPricesMin,
		orderBooks,
		greedIndex,
		news,
	))
	if err != nil {
		zlog.With(ctx).Warnw("Ask failed", "err", err)
		return nil, err
	}

	return decision, nil
}

func (d *dealUsecase) buy(ctx context.Context) (err error) {
	zlog.With(ctx).Infow("buy start")

	balance, err := d.upbitBankRepo.GetBalance(ctx)
	if err != nil {
		zlog.With(ctx).Warnw("Get balance failed", "err", err)
		return err
	}

	zlog.With(ctx).Infow("Got balance", "balance", balance)
	amount := balance.GetBuyAmount(d.feePercent, d.feeScale)
	if amount < d.minBuyAmount {
		zlog.With(ctx).Infow("No KRW balance")
		return nil
	}

	err = d.upbitBankRepo.Buy(ctx, amount)
	if err != nil {
		zlog.With(ctx).Warnw("Buy failed", "err", err)
		return err
	}

	return nil
}

func (d *dealUsecase) sell(ctx context.Context) (err error) {
	zlog.With(ctx).Infow("sell start")

	balance, err := d.upbitBankRepo.GetBitCoinBalance(ctx)
	if err != nil {
		zlog.With(ctx).Warnw("Get bitcoin balance failed", "err", err)
		return err
	}

	zlog.With(ctx).Infow("Got balance", "balance", balance)
	amount := balance.GetSellAmount()
	if amount == 0 {
		zlog.With(ctx).Infow("No bitcoin balance")
		return nil
	}

	err = d.upbitBankRepo.Sell(ctx, amount)
	if err != nil {
		zlog.With(ctx).Warnw("Sell failed", "err", err)
		return err
	}

	return nil
}
