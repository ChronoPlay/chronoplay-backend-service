package dto

import "github.com/ChronoPlay/chronoplay-backend-service/model"

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
}

type GetTransactionsRequest struct {
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
