package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
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
	coinBalance *BankBalance,
	totalDeposit float64,
	totalWithdrawal float64,
) *TransactionAggregate {
	transaction := NewTransaction(trResult)
	transactionSummary := NewTransactionSummary(
		trResult,
		totalDeposit,
		totalWithdrawal,
		bankBalance,
		coinBalance,
	)

	return &TransactionAggregate{
		Transaction:        *transaction,
		TransactionSummary: transactionSummary,
	}
}

// SetID ...
func (t *TransactionAggregate) SetID() {
	t.Transaction.SetID()
	t.TransactionSummary.SetID()
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

// String ...
func (t TransactionType) String() string {
	return toStr(transactionTypeStr, int(t))
}

// MarshalJSON ...
func (t *TransactionType) MarshalJSON() (data []byte, err error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON ...
func (t *TransactionType) UnmarshalJSON(data []byte) (err error) {
	*t = TransactionType(unmarshalJSON(data, transactionTypeMap, int(TransactionTypeNone)))
	return nil
}

// Transaction ...
type Transaction struct {
	ID         uint32
	BankID     uint32
	Identifier *string
	TrType     TransactionType
	Stock      string
	Percent    uint
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
		Percent:    trResult.Percent,
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

// SetID ...
func (t *Transaction) SetID() {
	t.ID = uuid.New().ID()
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
	coinBalance *BankBalance,
) *TransactionSummary {
	ts := &TransactionSummary{
		TotalDeposit:    totalDeposit,
		TotalWithdrawal: totalWithdrawal,
	}

	if bankBalance != nil {
		ts.TotalBankBalance = bankBalance.Balance
		ts.StockSummaries = append(ts.StockSummaries, NewTransactionStockSummary(bankBalance))
	}

	if coinBalance != nil {
		ts.TotalInvestBalance = coinBalance.Balance * coinBalance.AvgBuyPrice
		ts.StockSummaries = append(ts.StockSummaries, NewTransactionStockSummary(coinBalance))
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

// SetID ...
func (t *TransactionSummary) SetID() {
	t.ID = uuid.New().ID()
	t.StockSummaries.SetIDs()
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

// SetID ...
func (t *TransactionStockSummary) SetID() {
	t.ID = uuid.New().ID()
}

// TransactionStockSummaries ...
type TransactionStockSummaries []*TransactionStockSummary

// SetIDs ...
func (ts TransactionStockSummaries) SetIDs() {
	for _, v := range ts {
		v.SetID()
	}
}
