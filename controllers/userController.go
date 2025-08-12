package controller

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/ChronoPlay/chronoplay-backend-service/constants"
	"github.com/ChronoPlay/chronoplay-backend-service/mapper"
	"github.com/ChronoPlay/chronoplay-backend-service/model"
	service "github.com/ChronoPlay/chronoplay-backend-service/services"
)

type userController struct {
	userService service.UserService
}

type UserController interface {
	RegisterUser(*gin.Context)
	VerifyUser(*gin.Context)
	LoginUser(*gin.Context)
	GetUser(*gin.Context)
}

func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (ctl *userController) RegisterUser(c *gin.Context) {
	fmt.Print("request has reached here - userController")
	user, err := mapper.DecodeRegisterUserRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}

	ctx := c.Request.Context()
	err = ctl.userService.RegisterUser(ctx, user)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    "",
		Message: "User registered successfully. Checkout your mail to verify your emailId",
	})
}

func (ctl *userController) VerifyUser(c *gin.Context) {
	req, err := mapper.DecodeVerifyUserRequest(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	ctx := c.Request.Context()
	err = ctl.userService.VerifyUser(ctx, req)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    "",
		Message: "User successfully verified",
	})
}

func (ctl *userController) LoginUser(c *gin.Context) {
	req, err := mapper.DecodeLoginUser(c)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	ctx := c.Request.Context()
	resp, err := ctl.userService.LoginUser(ctx, req)
	log.Printf("Error type: %T, value: %#v\n", err, err)
	if err != nil {

		// Try logging the error type and value

		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    resp,
		Message: "User logged in successfully",
	})
}

func (ctl *userController) GetUser(c *gin.Context) {
	ctx := c.Request.Context()
	userId, _ := c.Get("UserID")
	user, err := ctl.userService.GetUser(ctx, model.User{
		UserId: userId.(uint32),
	})
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(200, constants.JsonResp{
		Data:    user,
		Message: "User fetched successfully",
	})
}
