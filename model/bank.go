package model

import (
	"encoding/json"
	"magmar/model/dao"
	"magmar/util"
	"math"
	"strconv"
)

// BankBalance ...
type BankBalance struct {
	Currency    string
	Balance     float64
	AvgBuyPrice float64
}

// NewBankBalance ...
func NewBankBalance(currency string, balance float64, avgBuyPrice float64) *BankBalance {
	return &BankBalance{
		Currency:    currency,
		Balance:     balance,
		AvgBuyPrice: avgBuyPrice,
	}
}

// GetBuyAmount can't be float64, calculate with uint64
func (b *BankBalance) GetBuyAmount(percent uint, feePercent uint, feeScale uint) (amount uint64) {
	price := uint64(b.Balance) * uint64(percent) / 100

	scale := util.Pow10(uint64(feeScale))
	feeAmountFloat := float64(price*uint64(feePercent)) / float64(scale)
	feeAmount := uint64(math.Ceil(feeAmountFloat))

	return price - feeAmount
}

// GetSellAmount ...
func (b *BankBalance) GetSellAmount(percent uint) (amount float64) {
	return b.Balance * float64(percent) / 100
}

// BankTransactionResult ...
type BankTransactionResult struct {
	Stock      string          `json:"stock"`
	TrType     TransactionType `json:"tr_type"`
	Price      *float64        `json:"price"`
	Volume     *float64        `json:"volume"`
	Identifier *string         `json:"identifier"`
	Reason     *string         `json:"reason"`
	Remark     *string         `json:"remark"`
}

// NewBankTransactionResultBuy ...
func NewBankTransactionResultBuy(result *dao.UpbitTransactionResult) *BankTransactionResult {
	bResult := &BankTransactionResult{
		TrType: TransactionTypeBuy,
		Stock:  string(result.Market),
	}

	if result.Price != nil {
		if price, err := strconv.ParseFloat(*result.Price, 64); err == nil {
			bResult.Price = &price
		}
	}

	if result.Volume != nil {
		if volume, err := strconv.ParseFloat(*result.Volume, 64); err == nil {
			bResult.Volume = &volume
		}
	}

	if result.Identifier != nil {
		bResult.Identifier = result.Identifier
	}

	return bResult
}

// NewBankTransactionResultBuyFailed ...
func NewBankTransactionResultBuyFailed(remark string) *BankTransactionResult {
	return &BankTransactionResult{
		Stock:  string(dao.UpbitStockBTC),
		TrType: TransactionTypeBuyFailed,
		Remark: &remark,
	}
}

// NewBankTransactionResultSell ...
func NewBankTransactionResultSell(result *dao.UpbitTransactionResult) *BankTransactionResult {
	bResult := &BankTransactionResult{
		Stock:  string(result.Market),
		TrType: TransactionTypeSell,
	}

	if result.Price != nil {
		if price, err := strconv.ParseFloat(*result.Price, 64); err == nil {
			bResult.Price = &price
		}
	}

	if result.Volume != nil {
		if volume, err := strconv.ParseFloat(*result.Volume, 64); err == nil {
			bResult.Volume = &volume
		}
	}

	if result.Identifier != nil {
		bResult.Identifier = result.Identifier
	}

	return bResult
}

// NewBankTransactionResultSellFailed ...
func NewBankTransactionResultSellFailed(remark string) *BankTransactionResult {
	return &BankTransactionResult{
		Stock:  string(StockNameUpbitBTC),
		TrType: TransactionTypeSellFailed,
		Remark: &remark,
	}
}

// NewBankTransactionResultHold ...
func NewBankTransactionResultHold() *BankTransactionResult {
	return &BankTransactionResult{
		Stock:  string(StockNameUpbitBTC),
		TrType: TransactionTypeHold,
	}
}

// NewBankTransactionResultDeposit ...
func NewBankTransactionResultDeposit(volume float64) *BankTransactionResult {
	return &BankTransactionResult{
		Stock:  string(StockNameBank),
		TrType: TransactionTypeDeposit,
		Volume: &volume,
		Price:  util.ToPtr(1.0),
	}
}

// NewBankTransactionResultWithdrawal ...
func NewBankTransactionResultWithdrawal(volume float64) *BankTransactionResult {
	return &BankTransactionResult{
		Stock:  string(StockNameBank),
		TrType: TransactionTypeWithdrawal,
		Volume: &volume,
		Price:  util.ToPtr(1.0),
	}
}

// String ...
func (b *BankTransactionResult) String() string {
	bytes, err := json.Marshal(b)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

// NewTotalDeposit ...
func (b *BankTransactionResult) NewTotalDeposit(totalDeposit float64) float64 {
	if b.TrType == TransactionTypeDeposit && b.Volume != nil {
		return *b.Volume + totalDeposit
	}
	return totalDeposit
}

// NewTotalWithdrawal ...
func (b *BankTransactionResult) NewTotalWithdrawal(totalWithdrawal float64) float64 {
	if b.TrType == TransactionTypeWithdrawal && b.Volume != nil {
		return *b.Volume + totalWithdrawal
	}
	return totalWithdrawal
}

// SetReason ...
func (b *BankTransactionResult) SetReason(reason string) {
	if reason == "" {
		return
	}

	b.Reason = util.ToPtr(reason)
}
