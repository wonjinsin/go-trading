package model

import (
	"time"
)

// BankID ...
const (
	BankMainID uint32 = 1
	BankSubID         = 2
)

// TransactionAggregate ...
type TransactionAggregate struct {
	Transaction
	TransactionSummary *TransactionSummary `gorm:"foreignKey:TransactionID"`
}

// NewTransactionAggregate ...
func NewTransactionAggregate(
	trResult *BankTransactionResult,
	bankBalance *BankBalance,
	bitCoinBalance *BankBalance,
	totalDeposit float64,
	totalWithdrawal float64,
) *TransactionAggregate {
	transaction := NewTransaction(trResult)
	transactionSummary := NewTransactionSummary(
		trResult,
		totalDeposit,
		totalWithdrawal,
		bankBalance,
		bitCoinBalance,
	)

	return &TransactionAggregate{
		Transaction:        *transaction,
		TransactionSummary: transactionSummary,
	}
}

// TransactionAggregates ...
type TransactionAggregates []*TransactionAggregate

// TransactionType ...
type TransactionType int32

// TransactionTypeConst ...
const (
	TransactionTypeNone TransactionType = iota
	TransactionTypeBuy
	TransactionTypeSell
	TransactionTypeBuyFailed
	TransactionTypeSellFailed
	TransactionTypeHold
	TransactionTypeDeposit
	TransactionTypeWithdrawal
)

var transactionTypeStr []string = []string{"None", "Buy", "Sell", "BuyFailed", "SellFailed", "Hold", "Deposit", "Withdrawal"}
var transactionTypeMap = map[string]int{
	"none":       int(TransactionTypeNone),
	"buy":        int(TransactionTypeBuy),
	"sell":       int(TransactionTypeSell),
	"buyfailed":  int(TransactionTypeBuyFailed),
	"sellfailed": int(TransactionTypeSellFailed),
	"hold":       int(TransactionTypeHold),
	"deposit":    int(TransactionTypeDeposit),
	"withdrawal": int(TransactionTypeWithdrawal),
}

// Transaction ...
type Transaction struct {
	ID         uint32
	BankID     uint32
	Identifier *string
	TrType     TransactionType
	Stock      string
	Volume     float64
	Price      float64
	Reason     *string
	Remark     *string
	CreatedAt  time.Time
}

// NewTransaction ...
func NewTransaction(trResult *BankTransactionResult) *Transaction {
	tr := &Transaction{
		BankID:     BankMainID,
		Identifier: trResult.Identifier,
		TrType:     trResult.TrType,
		Stock:      trResult.Stock,
		Remark:     trResult.Remark,
		CreatedAt:  time.Now(),
	}

	if trResult.Volume != nil {
		tr.Volume = *trResult.Volume
	}

	if trResult.Price != nil {
		tr.Price = *trResult.Price
	}

	if trResult.Reason != nil {
		tr.Reason = trResult.Reason
	}

	if trResult.Remark != nil {
		tr.Remark = trResult.Remark
	}

	return tr
}

// TableName ...
func (Transaction) TableName() string {
	return "transaction"
}

// TransactionSummary ...
type TransactionSummary struct {
	ID                 uint32
	TransactionID      uint32
	TotalDeposit       float64
	TotalWithdrawal    float64
	TotalBankBalance   float64
	TotalInvestBalance float64
	TotalProfit        float64
	StockSummaries     TransactionStockSummaries `gorm:"foreignKey:TransactionSummaryID"`
}

// NewTransactionSummary ...
func NewTransactionSummary(
	trResult *BankTransactionResult,
	totalDeposit float64,
	totalWithdrawal float64,
	bankBalance *BankBalance,
	bitCoinBalance *BankBalance,
) *TransactionSummary {
	ts := &TransactionSummary{
		TotalDeposit:    totalDeposit,
		TotalWithdrawal: totalWithdrawal,
	}

	if bankBalance != nil {
		ts.TotalBankBalance = bankBalance.Balance
		ts.StockSummaries = append(ts.StockSummaries, NewTransactionStockSummary(bankBalance))
	}

	if bitCoinBalance != nil {
		ts.TotalInvestBalance = bitCoinBalance.Balance * bitCoinBalance.AvgBuyPrice
		ts.StockSummaries = append(ts.StockSummaries, NewTransactionStockSummary(bitCoinBalance))
	}

	currentBalance := ts.TotalBankBalance + ts.TotalInvestBalance
	if totalDeposit != 0 {
		ts.TotalProfit = (currentBalance - totalDeposit + totalWithdrawal) / totalDeposit * 100
	}

	return ts
}

// TableName ...
func (TransactionSummary) TableName() string {
	return "transaction_summary"
}

// TransactionStockSummary ...
type TransactionStockSummary struct {
	ID                   uint32
	TransactionSummaryID uint32
	Stock                string
	Volume               float64
	Price                float64
}

// NewTransactionStockSummary ...
func NewTransactionStockSummary(balance *BankBalance) *TransactionStockSummary {
	return &TransactionStockSummary{
		Stock:  balance.Currency,
		Volume: balance.Balance,
		Price:  balance.AvgBuyPrice,
	}
}

// TableName ...
func (TransactionStockSummary) TableName() string {
	return "transaction_stock_summary"
}

// TransactionStockSummaries ...
type TransactionStockSummaries []*TransactionStockSummary
