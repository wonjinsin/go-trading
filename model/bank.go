package model

import "magmar/util"

// BankBalance ...
type BankBalance struct {
	Currency string
	Balance  uint64
}

// GetBuyAmount ...
func (b *BankBalance) GetBuyAmount(feePercent uint, feeScale uint) (amount uint64) {
	// (balance * feePercent + 99) / 100 for ceil
	feeAmount := (b.Balance*uint64(feePercent) + util.Pow10(uint64(feeScale)) - 1) / util.Pow10(uint64(feeScale))
	return b.Balance - feeAmount
}
