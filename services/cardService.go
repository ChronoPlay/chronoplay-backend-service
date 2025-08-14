package service

import (
	model "github.com/ChronoPlay/chronoplay-backend-service/model"
)

type CardService interface {
}

type cardService struct {
	cardRepo model.CardRepository
}

func NewCardService(cardRepo model.CardRepository) CardService {
	return &cardService{
		cardRepo: cardRepo,
	}
}
