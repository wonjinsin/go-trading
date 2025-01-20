package service

import (
	"context"
	"magmar/config"
	"magmar/model"
	"magmar/prompt"
	"magmar/repository"
	"magmar/util"

	"github.com/juju/errors"
)

type dealUsecase struct {
	conf            *config.ViperConfig
	feePercent      uint
	feeScale        uint
	minBuyAmount    uint64
	openAIqaRepo    repository.QaRepository
	upbitBankRepo   repository.StockBankRepository
	greedRepo       repository.GreedRepository
	newsRepo        repository.NewsRepository
	transactionRepo repository.TransactionRepository
	telegramMsgRepo repository.MsgRepository
}

// NewDealService ...
func NewDealService(conf *config.ViperConfig,
	qaRepo repository.QaRepository,
	upbitBankRepo repository.StockBankRepository,
	greedRepo repository.GreedRepository,
	newsRepo repository.NewsRepository,
	transactionRepo repository.TransactionRepository,
	telegramMsgRepo repository.MsgRepository) DealService {
	return &dealUsecase{
		conf:            conf,
		feePercent:      conf.GetUint(util.UpbitFeePercent),
		feeScale:        conf.GetUint(util.UpbitFeeScale),
		minBuyAmount:    conf.GetUint64(util.UpbitMinBuyAmount),
		openAIqaRepo:    qaRepo,
		upbitBankRepo:   upbitBankRepo,
		greedRepo:       greedRepo,
		newsRepo:        newsRepo,
		transactionRepo: transactionRepo,
		telegramMsgRepo: telegramMsgRepo,
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

	var trResult *model.BankTransactionResult
	switch decision.Decision {
	case model.DecisionStateBuy:
		trResult, err = d.buy(ctx, decision.Percent)
	case model.DecisionStateSell:
		trResult, err = d.sell(ctx, decision.Percent)
	case model.DecisionStateHold:
		trResult = model.NewBankTransactionResultHold()
	}
	trResult.SetReason(decision.Reason)

	if err != nil {
		zlog.With(ctx).Warnw("Decision handle failed", "err", err)
		return err
	}

	if _, err = d.saveTransaction(ctx, trResult); err != nil {
		zlog.With(ctx).Warnw("Save transaction failed", "err", err)
		return err
	}

	err = d.telegramMsgRepo.SendMessage(ctx, trResult.String())
	if err != nil {
		zlog.With(ctx).Warnw("Send message failed", "err", err)
	}

	zlog.With(ctx).Infow("Process done")

	return nil
}

func (d *dealUsecase) ask(ctx context.Context) (decision *model.Decision, err error) {
	marketPricesDay, err := d.upbitBankRepo.GetMarketPriceDataDay(ctx, model.StockNameUpbitBTC, 60)
	if err != nil {
		zlog.With(ctx).Warnw("Get market price data day failed", "err", err)
		return nil, err
	}

	marketPricesMin, err := d.upbitBankRepo.GetMarketPriceDataMin(ctx, model.StockNameUpbitBTC, 60)
	if err != nil {
		zlog.With(ctx).Warnw("Get market price data min failed", "err", err)
		return nil, err
	}

	orderBooks, err := d.upbitBankRepo.GetOrderBooks(ctx, model.StockNameUpbitBTC)
	if err != nil {
		zlog.With(ctx).Warnw("Get order book failed", "err", err)
		return nil, err
	}

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

func (d *dealUsecase) buy(ctx context.Context, percentage uint) (trResult *model.BankTransactionResult, err error) {
	zlog.With(ctx).Infow("buy start", "percentage", percentage)

	if percentage == 0 {
		zlog.With(ctx).Errorw("No percentage")
		return nil, errors.New("no percentage")
	}

	balance, err := d.upbitBankRepo.GetBalance(ctx)
	if err != nil {
		zlog.With(ctx).Warnw("Get balance failed", "err", err)
		return nil, err
	}

	zlog.With(ctx).Infow("Got balance", "balance", balance)
	amount := balance.GetBuyAmount(percentage, d.feePercent, d.feeScale)
	if amount < d.minBuyAmount {
		zlog.With(ctx).Infow("No KRW balance")
		return model.NewBankTransactionResultBuyFailed(util.RemarkBankTransactionResultBuyFailedMinAmount), nil
	}

	trResult, err = d.upbitBankRepo.Buy(ctx, amount)
	if err != nil {
		zlog.With(ctx).Warnw("Buy failed", "err", err)
		return nil, err
	}

	return trResult, nil
}

func (d *dealUsecase) sell(ctx context.Context, percentage uint) (trResult *model.BankTransactionResult, err error) {
	zlog.With(ctx).Infow("sell start", "percentage", percentage)

	if percentage == 0 {
		zlog.With(ctx).Errorw("No percentage")
		return nil, errors.New("No percentage")
	}

	balance, err := d.upbitBankRepo.GetBitCoinBalance(ctx)
	if err != nil {
		if errors.Is(err, errors.NotFound) {
			zlog.With(ctx).Infow("No bitcoin balance")
			return model.NewBankTransactionResultSellFailed(util.RemarkBankTransactionResultSellFailedMinAmount), nil
		}
		zlog.With(ctx).Warnw("Get bitcoin balance failed", "err", err)
		return nil, err
	}

	zlog.With(ctx).Infow("Got balance", "balance", balance)
	amount := balance.GetSellAmount(percentage)
	if amount == 0 {
		zlog.With(ctx).Infow("No bitcoin balance")
		return model.NewBankTransactionResultSellFailed(util.RemarkBankTransactionResultSellFailedMinAmount), nil
	}

	trResult, err = d.upbitBankRepo.Sell(ctx, amount)
	if err != nil {
		zlog.With(ctx).Warnw("Sell failed", "err", err)
		return nil, err
	}

	return trResult, nil
}

func (d *dealUsecase) saveTransaction(ctx context.Context, trResult *model.BankTransactionResult) (transaction *model.TransactionAggregate, err error) {
	zlog.With(ctx).Infow("saveTransaction start", "trResult", trResult)

	bankBalance, err := d.upbitBankRepo.GetBalance(ctx)
	if err != nil {
		zlog.With(ctx).Warnw("Get balance failed", "err", err)
		return nil, err
	}

	bitCoinBalance, err := d.upbitBankRepo.GetBitCoinBalance(ctx)
	if err != nil && !errors.Is(err, errors.NotFound) {
		zlog.With(ctx).Warnw("Get bitcoin balance failed", "err", err)
		return nil, err
	}

	totalDeposit, err := d.transactionRepo.GetTotalDeposit(ctx)
	if err != nil {
		zlog.With(ctx).Warnw("Get total deposit failed", "err", err)
		return nil, err
	}

	totalWithdrawal, err := d.transactionRepo.GetTotalWithdrawal(ctx)
	if err != nil {
		zlog.With(ctx).Warnw("Get total withdrawal failed", "err", err)
		return nil, err
	}

	transaction = model.NewTransactionAggregate(trResult, bankBalance, bitCoinBalance, totalDeposit, totalWithdrawal)

	transaction, err = d.transactionRepo.NewTransaction(ctx, transaction)
	if err != nil {
		zlog.With(ctx).Warnw("Save transaction failed", "err", err)
		return nil, err
	}

	return transaction, nil
}
