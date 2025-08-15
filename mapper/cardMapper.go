package mapper

import (
	"log"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	"github.com/ChronoPlay/chronoplay-backend-service/model"
	"github.com/gin-gonic/gin"
)

func DecodeAddCardRequest(c *gin.Context) (dto.AddCardRequest, *helpers.CustomError) {
	var req dto.AddCardRequest
	if err := c.ShouldBind(&req); err != nil {
		return dto.AddCardRequest{}, helpers.BadRequest("Invalid request body" + err.Error())
	}
	log.Println("Decoded AddCardRequest:", req)
	userId, _ := c.Get("UserID")
	req.UserId = userId.(uint32)
	file, err := c.FormFile("image")
	if err != nil {
		return dto.AddCardRequest{}, helpers.BadRequest("Failed to get image from request")
	}
	req.Image = file
	return req, nil
}

func DecodeGetCardRequest(c *gin.Context) (dto.GetCardRequest, *helpers.CustomError) {
	var req dto.GetCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return dto.GetCardRequest{}, helpers.BadRequest("Invalid request body" + err.Error())
	}
	return req, nil
}

func EncodeGetCardResponse(req *model.Card) (res dto.GetCardResponse) {
	res.Name = req.Name
	res.Description = req.Description
	res.Number = req.Number
	res.Total = req.Total
	res.Occupied = req.Occupied
	res.ImageUrl = req.ImageUrl
	return res
}
