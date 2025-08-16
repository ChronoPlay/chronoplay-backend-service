package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	GetUserById(*gin.Context)
	AddFriend(*gin.Context)
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
	req, err := mapper.DecodeLoginUserRequest(c)
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
func (ctl *userController) GetUserById(c *gin.Context) {
	ctx := c.Request.Context()
	userId := c.Query("user_id")
	id64, perr := strconv.ParseUint(userId, 10, 32)
	if perr != nil {
		c.JSON(http.StatusBadRequest, constants.JsonResp{
			Message: "Error while parsing user_id" + perr.Error(),
		})
		return
	}
	user, err := ctl.userService.GetUser(ctx, model.User{
		UserId: uint32(id64),
	})
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusBadRequest, constants.JsonResp{
			Message: "user not found by the user_id",
		})
		return
	}
	data := mapper.EncodeGetUserResponse(user)

	c.JSON(200, constants.JsonResp{
		Data:    data,
		Message: "User details fetched successfully",
	})
}
func (ctl *userController) AddFriend(c *gin.Context) {
	ctx := c.Request.Context()
	CurUserId, _ := c.Get("UserID")
	FriendUserId := c.Query("user_id")
	fid, perr := strconv.ParseInt(FriendUserId, 10, 32)
	if perr != nil {
		c.JSON(http.StatusBadRequest, constants.JsonResp{
			Message: "Error while parsing user_id" + perr.Error(),
		})
		return
	}
	nerr := ctl.userService.AddFriend(ctx, CurUserId.(uint32), uint32(fid))
	if nerr != nil {
		c.JSON(int(nerr.Code), constants.JsonResp{
			Message: nerr.Message,
		})
		return
	}
	c.JSON(http.StatusOK, constants.JsonResp{
		Message: "friend added succesfully",
	})

}
