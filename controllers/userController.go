package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ChronoPlay/chronoplay-backend-service/constants"
	"github.com/ChronoPlay/chronoplay-backend-service/mapper"
	model "github.com/ChronoPlay/chronoplay-backend-service/models"
	service "github.com/ChronoPlay/chronoplay-backend-service/services"
)

type userController struct {
	userService service.UserService
}

type UserController interface {
	GetUser(*gin.Context)
	RegisterUser(*gin.Context)
}

func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (ctl *userController) GetUser(c *gin.Context) {
	c.JSON(200, model.User{
		Name: "Sparsh",
		Cash: 20020103,
	})
}

func (ctl *userController) RegisterUser(c *gin.Context) {
	user, err := mapper.DecodeRegisterUserRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Messgae: err.Message,
		})
	}

	ctx := c.Request.Context()
	err = ctl.userService.RegisterUser(ctx, user)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Messgae: err.Message,
		})
	}
	c.JSON(200, constants.JsonResp{
		Data: "",
	})
}
