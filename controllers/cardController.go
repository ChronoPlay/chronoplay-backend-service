package controller

import service "github.com/ChronoPlay/chronoplay-backend-service/services"

type cardController struct {
	cardService service.CardService
}

type CardController interface {
}

func NewCardController(cardService service.CardService) CardController {
	return &cardController{
		cardService: cardService,
	}
}
