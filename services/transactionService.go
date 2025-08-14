package service

import model "github.com/ChronoPlay/chronoplay-backend-service/model"

type TransactionService interface {
}

type transactionService struct {
	cardTransactionRepo model.CardTransactionRepository
	cashTransactionRepo model.CashTransactionRepository
}

func NewTransactionService(cardTransactionRepo model.CardTransactionRepository, cashTransactionRepo model.CashTransactionRepository) TransactionService {
	return &transactionService{
		cardTransactionRepo: cardTransactionRepo,
		cashTransactionRepo: cashTransactionRepo,
	}
}
