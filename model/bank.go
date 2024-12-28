package model

import "magmar/util"

// BankBalance ...
type BankBalance struct {
	Currency string
	Balance  float64
}

// GetBuyAmount can't be float64, calculate with uint64
func (b *BankBalance) GetBuyAmount(feePercent uint, feeScale uint) (amount uint64) {
	// (balance * feePercent + 99) / 100 for ceil
	feeAmount := (uint64(b.Balance)*uint64(feePercent) + util.Pow10(uint64(feeScale)) - 1) / util.Pow10(uint64(feeScale))
	return uint64(b.Balance) - feeAmount
}

// GetSellAmount ...
func (b *BankBalance) GetSellAmount() (amount float64) {
	return b.Balance
}
