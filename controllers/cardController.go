package controller

import (
	service "github.com/ChronoPlay/chronoplay-backend-service/services"
	"github.com/gin-gonic/gin"
)

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

func (cardCtrl *cardController) AddCard(c *gin.Context) {

}
