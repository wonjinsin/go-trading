package model

import (
	"magmar/util"
	"math"
)

// BankBalance ...
type BankBalance struct {
	Currency string
	Balance  float64
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
