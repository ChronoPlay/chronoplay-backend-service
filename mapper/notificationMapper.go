package mapper

import (
	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	"github.com/gin-gonic/gin"
)

func DecodeGetNotificationsRequest(c *gin.Context) (req dto.GetNotificationsRequest, err *helpers.CustomError) {
	userId, _ := c.Get("UserID")
	req.UserId = userId.(uint32)
	return req, nil
}

func DecodeMarkNotificationsAsReadRequest(c *gin.Context) (req dto.MarkNotificationsAsReadRequest, err *helpers.CustomError) {
	userId, _ := c.Get("UserID")
	req.UserId = userId.(uint32)
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, helpers.BadRequest("Invalid request payload")
	}
	return req, nil
}
