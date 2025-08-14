package controller

import service "github.com/ChronoPlay/chronoplay-backend-service/services"

type transactionController struct {
	transactionService service.TransactionService
}

type TransactionController interface {
}

func NewTransactionController(transactionService service.TransactionService) TransactionController {
	return &transactionController{
		transactionService: transactionService,
	}
}
