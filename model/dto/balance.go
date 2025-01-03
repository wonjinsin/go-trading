package dto

// Deposit ...
type Deposit struct {
	Price float64 `json:"price"`
}

// Withdrawal ...
type Withdrawal struct {
	Deposit
}
