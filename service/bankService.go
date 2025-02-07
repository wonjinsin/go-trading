package service

import (
	"context"
	"magmar/model"
	"magmar/model/dto"
	"magmar/repository"
	"magmar/util"
)

// bankUsecase ...
type bankUsecase struct {
	upbitBankRepo   repository.StockBankRepository
	transactionRepo repository.TransactionRepository
}

// NewBankService ...
func NewBankService(upbitBank repository.StockBankRepository, transaction repository.TransactionRepository) BankService {
	return &bankUsecase{
		upbitBankRepo:   upbitBank,
		transactionRepo: transaction,
	}
}

// Deposit ...
func (b *bankUsecase) Deposit(ctx context.Context, deposit dto.Deposit) (*model.TransactionAggregate, error) {
	zlog.With(ctx).Infow(util.LogSvc, "deposit", deposit)
	trResult := model.NewBankTransactionResultDeposit(deposit.Price)
	return b.saveTransaction(ctx, trResult)
}

// Withdrawal ...
func (b *bankUsecase) Withdrawal(ctx context.Context, withdrawal dto.Withdrawal) (*model.TransactionAggregate, error) {
	zlog.With(ctx).Infow(util.LogSvc, "withdrawal", withdrawal)
	trResult := model.NewBankTransactionResultWithdrawal(withdrawal.Price)
	return b.saveTransaction(ctx, trResult)
}

func (b *bankUsecase) saveTransaction(ctx context.Context, trResult *model.BankTransactionResult) (transaction *model.TransactionAggregate, err error) {
	zlog.With(ctx).Infow("saveTransaction start", "trResult", trResult)
	// pass because not use db
	return nil, nil

	// bankBalance, err := b.upbitBankRepo.GetBalance(ctx)
	// if err != nil {
	// 	zlog.With(ctx).Warnw("Get balance failed", "err", err)
	// 	return nil, err
	// }

	// bitCoinBalance, err := b.upbitBankRepo.GetCoinBalance(ctx, util.CoinMap[util.BitCoin][util.UpbitCurrency])
	// if err != nil && !errors.Is(err, errors.NotFound) {
	// 	zlog.With(ctx).Warnw("Get bitcoin balance failed", "err", err)
	// 	return nil, err
	// }

	// totalDeposit, err := b.transactionRepo.GetTotalDeposit(ctx)
	// totalDeposit = trResult.NewTotalDeposit(totalDeposit)
	// if err != nil {
	// 	zlog.With(ctx).Warnw("Get total deposit failed", "err", err)
	// 	return nil, err
	// }

	// totalWithdrawal, err := b.transactionRepo.GetTotalWithdrawal(ctx)
	// totalWithdrawal = trResult.NewTotalWithdrawal(totalWithdrawal)
	// if err != nil {
	// 	zlog.With(ctx).Warnw("Get total withdrawal failed", "err", err)
	// 	return nil, err
	// }

	// transaction = model.NewTransactionAggregate(trResult, bankBalance, bitCoinBalance, totalDeposit, totalWithdrawal)

	// transaction, err = b.transactionRepo.NewTransaction(ctx, transaction)
	// if err != nil {
	// 	zlog.With(ctx).Warnw("Save transaction failed", "err", err)
	// 	return nil, err
	// }

	// return transaction, nil
}
