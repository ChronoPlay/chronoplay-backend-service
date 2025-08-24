package controller

import (
	"github.com/ChronoPlay/chronoplay-backend-service/constants"
	"github.com/ChronoPlay/chronoplay-backend-service/mapper"
	service "github.com/ChronoPlay/chronoplay-backend-service/services"
	"github.com/gin-gonic/gin"
)

type NotificationController interface {
	GetNotifications(*gin.Context)
	MarkAsRead(*gin.Context)
}

type notificationController struct {
	notificationService service.NotificationService
}

func NewNotificationController(notificationService service.NotificationService) NotificationController {
	return &notificationController{
		notificationService: notificationService,
	}
}

func (ctl *notificationController) GetNotifications(c *gin.Context) {
	req, err := mapper.DecodeGetNotificationsRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	ctx := c.Request.Context()
	notifications, err := ctl.notificationService.GetNotifications(ctx, req)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    notifications,
		Message: "Notifications fetched successfully.",
	})
}

func (ctl *notificationController) MarkAsRead(c *gin.Context) {
	req, err := mapper.DecodeMarkNotificationsAsReadRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	ctx := c.Request.Context()
	err = ctl.notificationService.MarkNotificationsAsRead(ctx, req)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    "",
		Message: "Notifications marked as read successfully.",
	})
}
