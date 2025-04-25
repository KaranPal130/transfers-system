package models

import "github.com/shopspring/decimal"

// Account represents a financial account
type Account struct {
	AccountID int64           `json:"account_id"`
	Balance   decimal.Decimal `json:"balance"`
}

// AccountCreateRequest represents the request to create a new account
type AccountCreateRequest struct {
	AccountID      int64  `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}
