package model

import "testing"

func TestBankBalance_GetBuyAmount(t *testing.T) {
	tests := []struct {
		name           string
		bankBalance    BankBalance
		percent        uint
		feePercent     uint
		feeScale       uint
		expectedAmount uint64
	}{
		{
			name: "100% with 0.005% fee",
			bankBalance: BankBalance{
				Balance: 10000.0,
			},
			percent:        100,
			feePercent:     5,
			feeScale:       4,
			expectedAmount: 9995,
		},
		{
			name: "20% with 0.005% fee",
			bankBalance: BankBalance{
				Balance: 10000.0,
			},
			percent:        20,
			feePercent:     5,
			feeScale:       4,
			expectedAmount: 1999,
		},
		{
			name: "100% with no fee",
			bankBalance: BankBalance{
				Balance: 1000.0,
			},
			percent:        100,
			feePercent:     0,
			feeScale:       8,
			expectedAmount: 1000,
		},
		{
			name: "50% with 0.1% fee",
			bankBalance: BankBalance{
				Balance: 1000.0,
			},
			percent:        50,
			feePercent:     1,
			feeScale:       3,
			expectedAmount: 499, // (1000 * 50% = 500) - (500 * 0.001 = 0.5 rounded up to 1) = 499
		},
		{
			name: "25% with 0.2% fee",
			bankBalance: BankBalance{
				Balance: 2000.0,
			},
			percent:        25,
			feePercent:     2,
			feeScale:       3,
			expectedAmount: 499, // (2000 * 25% = 500) - (500 * 0.002 = 1) = 499
		},
		{
			name: "zero balance",
			bankBalance: BankBalance{
				Balance: 0.0,
			},
			percent:        100,
			feePercent:     1,
			feeScale:       3,
			expectedAmount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount := tt.bankBalance.GetBuyAmount(tt.percent, tt.feePercent, tt.feeScale)
			if amount != tt.expectedAmount {
				t.Errorf("GetBuyAmount(%v, %v, %v) = %v, want %v",
					tt.percent, tt.feePercent, tt.feeScale, amount, tt.expectedAmount)
			}
		})
	}
}
