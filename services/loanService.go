package service

import model "github.com/ChronoPlay/chronoplay-backend-service/model"

type LoanService interface {
}

type loanService struct {
	loanRepo model.LoanRepository
}

func NewLoanService(loanRepo model.LoanRepository) LoanService {
	return &loanService{
		loanRepo: loanRepo,
	}
}
