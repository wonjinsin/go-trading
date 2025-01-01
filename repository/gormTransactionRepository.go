package repository

import (
	"context"
	"magmar/model"
	"magmar/util"

	"gorm.io/gorm"
)

type gormTransactionRepository struct {
	conn *gorm.DB
}

// NewGormTransactionRepository ...
func NewGormTransactionRepository(conn *gorm.DB) TransactionRepository {
	return &gormTransactionRepository{conn: conn}
}

// NewTransaction ...
func (g *gormTransactionRepository) NewTransaction(ctx context.Context, transaction *model.TransactionAggregate) (*model.TransactionAggregate, error) {
	zlog.With(ctx).Infow(util.LogRepo, "transaction", transaction)

	scope := g.conn.WithContext(ctx).Begin()
	scope = scope.Create(&transaction)
	if err := scope.Error; err != nil {
		scope.Rollback()
		zlog.With(ctx).Errorw("NewTransaction Error", "err", err)
		return nil, err
	}

	if err := scope.Commit().Error; err != nil {
		zlog.With(ctx).Errorw("NewTransaction Commit Error", "err", err)
		return nil, err
	}

	return transaction, nil
}

// GetTransactions ...
func (g *gormTransactionRepository) GetTransactions(ctx context.Context) (transactions model.TransactionAggregates, err error) {
	zlog.With(ctx).Infow(util.LogRepo)

	scope := g.conn.WithContext(ctx)
	scope = scope.Find(&transactions)
	if err = scope.Error; err != nil {
		zlog.With(ctx).Errorw("GetTransactions Error", "err", err)
		return nil, err
	}

	return transactions, nil
}

// GetTotalDeposit ...
func (g *gormTransactionRepository) GetTotalDeposit(ctx context.Context) (float64, error) {
	zlog.With(ctx).Infow(util.LogRepo)

	var total float64

	scope := g.conn.WithContext(ctx)
	scope = scope.Model(&model.Transaction{}).
		Where("tr_type = ?", model.TransactionTypeDeposit).
		Select("COALESCE(SUM(volume), 0) as total").
		Scan(&total)
	if err := scope.Error; err != nil {
		zlog.With(ctx).Errorw("GetTotalDeposit Error", "err", err)
		return 0, err
	}

	return total, nil
}

// GetTotalWithdrawal ...
func (g *gormTransactionRepository) GetTotalWithdrawal(ctx context.Context) (float64, error) {
	zlog.With(ctx).Infow(util.LogRepo)

	var total float64

	scope := g.conn.WithContext(ctx)
	scope = scope.Model(&model.Transaction{}).
		Where("tr_type = ?", model.TransactionTypeWithdrawal).
		Select("COALESCE(SUM(volume), 0) as total").
		Scan(&total)
	if err := scope.Error; err != nil {
		zlog.With(ctx).Errorw("GetTotalWithdrawal Error", "err", err)
		return 0, err
	}

	return total, nil
}
