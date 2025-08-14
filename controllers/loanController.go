package controller

import service "github.com/ChronoPlay/chronoplay-backend-service/services"

type loanController struct {
	loanService service.LoanService
}

type LoanController interface {
}

func NewLoanController(loanService service.LoanService) LoanController {
	return &loanController{
		loanService: loanService,
	}
}
