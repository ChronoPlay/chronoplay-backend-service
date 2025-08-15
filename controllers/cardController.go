package controller

import (
	"net/http"

	"github.com/ChronoPlay/chronoplay-backend-service/constants"
	"github.com/ChronoPlay/chronoplay-backend-service/mapper"
	service "github.com/ChronoPlay/chronoplay-backend-service/services"
	"github.com/gin-gonic/gin"
)

type cardController struct {
	cardService service.CardService
}

type CardController interface {
	AddCard(c *gin.Context)
	GetCard(c *gin.Context)
}

func NewCardController(cardService service.CardService) CardController {
	return &cardController{
		cardService: cardService,
	}
}

func (cardCtrl *cardController) AddCard(c *gin.Context) {
	req, err := mapper.DecodeAddCardRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	ctx := c.Request.Context()
	err = cardCtrl.cardService.AddCard(ctx, req)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    "",
		Message: "Card added successfully",
	})
}

func (cardCtrl *cardController) GetCard(c *gin.Context) {
	req, err := mapper.DecodeGetCardRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	ctx := c.Request.Context()
	card, err := cardCtrl.cardService.GetCard(ctx, req)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	if card == nil {
		c.JSON(http.StatusBadRequest, constants.JsonResp{
			Message: "Card Not Found",
		})
		return
	}
	cardRes := mapper.EncodeGetCardResponse(card)
	c.JSON(200, constants.JsonResp{
		Data:    cardRes,
		Message: "Card Name Found successfully",
	})
}
