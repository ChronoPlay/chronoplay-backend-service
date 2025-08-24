package mapper

import (
	"strconv"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	"github.com/gin-gonic/gin"
)

func DecodeTransferCashRequest(c *gin.Context) (dto.TransferCashRequest, *helpers.CustomError) {
	var req dto.TransferCashRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, helpers.BadRequest("Invalid request body" + err.Error())
	}
	userId, _ := c.Get("UserID")
	req.UserId = userId.(uint32)
	return req, nil
}

func DecodeTransferCardsRequest(c *gin.Context) (dto.TransferCardRequest, *helpers.CustomError) {
	var req dto.TransferCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, helpers.BadRequest("Invalid request body" + err.Error())
	}
	userId, _ := c.Get("UserID")
	req.UserId = userId.(uint32)
	return req, nil
}

func DecodeGetTransactionsRequest(c *gin.Context) (dto.GetTransactionsRequest, *helpers.CustomError) {
	var req dto.GetTransactionsRequest
	userId, _ := c.Get("UserID")
	req.UserId = userId.(uint32)
	transactionGuid, exists := c.GetQuery("transaction_guid")
	if !exists {
		return req, helpers.BadRequest("Missing transaction_guid in query parameters")
	}
	transactionGuidUint, perr := strconv.ParseUint(transactionGuid, 10, 32)
	if perr != nil {
		return req, helpers.BadRequest("Invalid transaction_guid: " + perr.Error())
	}
	req.TransactionGuids = []uint32{uint32(transactionGuidUint)}
	return req, nil
}

func DecodeExchangeRequest(c *gin.Context) (dto.ExchangeRequest, *helpers.CustomError) {
	var req dto.ExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, helpers.BadRequest("Invalid request body" + err.Error())
	}
	userId, _ := c.Get("UserID")
	req.UserId = userId.(uint32)
	req.GivenBy = userId.(uint32)
	return req, nil
}

func DecodeGetPossibleExchangeRequest(c *gin.Context) (req dto.GetPossibleExchangeRequest, err *helpers.CustomError) {
	userId, _ := c.Get("UserID")
	req.UserId = userId.(uint32)
	traderId, exists := c.GetQuery("user_id")
	if !exists {
		return req, helpers.BadRequest("Missing user_id in query parameters")
	}
	traderIdUint, perr := strconv.ParseUint(traderId, 10, 32)
	if perr != nil {
		return req, helpers.BadRequest("Invalid trader_id: " + perr.Error())
	}
	req.TraderId = uint32(traderIdUint)
	return req, nil
}

func DecodeExecuteExchangeRequest(c *gin.Context) (dto.ExecuteExchangeRequest, *helpers.CustomError) {
	var req dto.ExecuteExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, helpers.BadRequest("Invalid request body" + err.Error())
	}
	userId, _ := c.Get("UserID")
	req.UserId = userId.(uint32)
	return req, nil
}
