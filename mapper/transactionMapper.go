package mapper

import (
	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	"github.com/gin-gonic/gin"
)

func DecodeTransferCashRequest(c *gin.Context) (dto.TransferCashRequest, *helpers.CustomError) {
	var req dto.TransferCashRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, helpers.BadRequest("Invalid request body")
	}
	userId, _ := c.Get("UserID")
	req.UserId = userId.(uint32)
	return req, nil
}

func DecodeTransferCardsRequest(c *gin.Context) (dto.TransferCardRequest, *helpers.CustomError) {
	var req dto.TransferCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, helpers.BadRequest("Invalid request body")
	}
	userId, _ := c.Get("UserID")
	req.UserId = userId.(uint32)
	return req, nil
}
