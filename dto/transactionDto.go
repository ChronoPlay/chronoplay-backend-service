package dto

import (
	"time"

	"github.com/ChronoPlay/chronoplay-backend-service/model"
)

type Transaction struct {
	TransactionGuid uint32    `json:"transaction_guid"`
	CardsRecieved   []Card    `json:"cards_recieved"`
	CardsSent       []Card    `json:"cards_sent"`
	CashSent        float32   `json:"cash_sent"`
	CashRecieved    float32   `json:"cash_recieved"`
	TransactionWith uint32    `json:"transaction_with"`
	Time            time.Time `json:"time"`
	Status          string    `json:"status"`
}

type Card struct {
	CardNumber string `json:"card_number"`
	Amount     uint32 `json:"amount"`
}

type TransferCashRequest struct {
	Amount   float32 `json:"amount"`
	GivenBy  uint32  `json:"given_by"`
	GivenTo  uint32  `json:"given_to"`
	Status   string  `json:"status"`
	UserId   uint32  `json:"user_id"`
	UserType string  `json:"user_type"`
}

type TransferCardRequest struct {
	Cards    []TransferCard `json:"cards"`
	GivenBy  uint32         `json:"given_by"`
	GivenTo  uint32         `json:"given_to"`
	Status   string         `json:"status"`
	UserId   uint32         `json:"user_id"`
	UserType string         `json:"user_type"`
}

type TransferCard struct {
	CardNumber string `json:"card_number"`
	Amount     uint32 `json:"amount"`
}

type ExchangeRequest struct {
	GivenBy       uint32  `json:"given_by"`
	GivenTo       uint32  `json:"given_to"`
	CashSent      float32 `json:"cash_sent"`
	CashRecieved  float32 `json:"cash_recieved"`
	CardsSent     []Card  `json:"cards_sent"`
	CardsRecieved []Card  `json:"cards_recieved"`
	UserId        uint32  `json:"user_id"`
	UserType      string  `json:"user_type"`
}

type GetTransactionsRequest struct {
	UserId           uint32   `json:"user_id"`
	TransactionGuids []uint32 `json:"transaction_guids"`
}

type GetTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type IsCashTransactionPossibleRequest struct {
	GivenBy uint32
	User    model.User
	Amount  float32
}

type IsCardTransactionPossibleRequest struct {
	GivenBy            uint32
	User               model.User
	CardsToTransferMap map[string]uint32
	CardsOccupiedMap   map[string]uint32
	CardsMap           map[string]model.Card
}

type IsExhangePossibleRequest struct {
	GivenByUser   model.User
	GivenToUser   model.User
	CashSent      float32
	CashRecieved  float32
	CardsSent     []Card
	CardsRecieved []Card
}

type GetPossibleExchangeRequest struct {
	UserId   uint32 `json:"user_id"`
	TraderId uint32 `json:"trader_id"`
}

type GetPossibleExchangeResponse struct {
	YourCash    float32        `json:"yourCash"`
	YourCards   []CardResponse `json:"yourCards"`
	TraderCash  float32        `json:"traderCash"`
	TraderCards []CardResponse `json:"traderCards"`
}

type ExecuteExchangeRequest struct {
	TransactionGuid uint32 `json:"transaction_guid"`
	UserId          uint32 `json:"user_id"`
	IsAccepted      bool   `json:"is_accepted"`
}

type IsvalidTransactionConfirmerRequest struct {
	UserId    uint32
	CreatedBy uint32
	GivenBy   uint32
	GivenTo   uint32
}
