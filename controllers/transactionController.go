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
	GetPossibleExchange(*gin.Context)
	ExecuteExchange(*gin.Context)
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
	req, err := mapper.DecodeExchangeRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	ctx := c.Request.Context()
	err = ctl.transactionService.Exchange(ctx, req)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    "",
		Message: "Exchange request created successfully. Please wait for the other user to accept the request.",
	})
}

func (ctl *transactionController) GetTransactions(c *gin.Context) {
	req, err := mapper.DecodeGetTransactionsRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	ctx := c.Request.Context()
	resp, err := ctl.transactionService.GetTransactions(ctx, req)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    resp.Transactions,
		Message: "Transactions fetched successfully",
	})
}

func (ctl *transactionController) GetPossibleExchange(c *gin.Context) {
	req, err := mapper.DecodeGetPossibleExchangeRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	ctx := c.Request.Context()
	resp, err := ctl.transactionService.GetPossibleExchange(ctx, req)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    resp,
		Message: "Possible exchanges fetched successfully",
	})
}

func (ctl *transactionController) ExecuteExchange(c *gin.Context) {
	req, err := mapper.DecodeExecuteExchangeRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	ctx := c.Request.Context()
	err = ctl.transactionService.ExecuteExchange(ctx, req)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    "",
		Message: "Exchange confirmed successfully",
	})
}
