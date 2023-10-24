package web

import "time"

type PaymentSuccess struct {
	Updated                  time.Time `json:"updated" binding:"required"`
	Created                  time.Time `json:"created" binding:"required"`
	PaymentID                string    `json:"payment_id" binding:"required"`
	CallbackVirtualAccountID string    `json:"callback_virtual_account_id" binding:"required"`
	OwnerID                  string    `json:"owner_id" binding:"required"`
	ExternalID               string    `json:"external_id" binding:"required"`
	AccountNumber            string    `json:"account_number" binding:"required"`
	BankCode                 string    `json:"bank_code" binding:"required"`
	Amount                   int       `json:"amount" binding:"required"`
	TransactionTimestamp     time.Time `json:"transaction_timestamp" binding:"required"`
	MerchantCode             string    `json:"merchant_code" binding:"required"`
	ID                       string    `json:"id" binding:"required"`
}
