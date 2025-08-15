package mapper

import (
	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	"github.com/gin-gonic/gin"
)

func DecodeAddCardRequest(c *gin.Context) (dto.AddCardRequest, *helpers.CustomError) {
	var req dto.AddCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return dto.AddCardRequest{}, helpers.BadRequest("Invalid request body")
	}
	userId, _ := c.Get("UserID")
	req.UserId = userId.(uint32)
	userType, _ := c.Get("UserType")
	req.UserType = userType.(string)
	return req, nil
}
