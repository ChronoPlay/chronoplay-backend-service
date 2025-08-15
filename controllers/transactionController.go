package controller

import (
	"github.com/ChronoPlay/chronoplay-backend-service/constants"
	"github.com/ChronoPlay/chronoplay-backend-service/mapper"
	service "github.com/ChronoPlay/chronoplay-backend-service/services"
	"github.com/gin-gonic/gin"
)

type transactionController struct {
	transactionService service.TransactionService
}

type TransactionController interface {
	Transfercash(*gin.Context)
	Transfercards(*gin.Context)
	Exchange(*gin.Context)
	GetTransactions(*gin.Context)
	GiveCards(*gin.Context)
}

func NewTransactionController(transactionService service.TransactionService) TransactionController {
	return &transactionController{
		transactionService: transactionService,
	}
}

func (ctl *transactionController) Transfercash(c *gin.Context) {
	req, err := mapper.DecodeTransferCashRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	ctx := c.Request.Context()
	err = ctl.transactionService.TransferCash(ctx, req)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    "",
		Message: "Cash transferred successfully",
	})
}

func (ctl *transactionController) Transfercards(c *gin.Context) {
	req, err := mapper.DecodeTransferCardsRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	ctx := c.Request.Context()
	err = ctl.transactionService.TransferCards(ctx, req)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    "",
		Message: "Cards transferred successfully",
	})

}

func (ctl *transactionController) Exchange(c *gin.Context) {
}

func (ctl *transactionController) GetTransactions(c *gin.Context) {
}

func (ctl *transactionController) GiveCards(c *gin.Context) {
	req, err := mapper.DecodeTransferCardsRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	ctx := c.Request.Context()
	err = ctl.transactionService.GiveCards(ctx, req)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    "",
		Message: "Cards given successfully",
	})
}
