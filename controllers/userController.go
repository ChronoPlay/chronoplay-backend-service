package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ChronoPlay/chronoplay-backend-service/constants"
	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/mapper"
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
	GetFriends(c *gin.Context)
	RemoveFriend(c *gin.Context)
	ActivateAllUsers(c *gin.Context)
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
	user, err := ctl.userService.GetUser(ctx, dto.GetUserRequest{
		UserID: userId.(uint32),
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
	user, err := ctl.userService.GetUser(ctx, dto.GetUserRequest{
		UserID: uint32(id64),
	})
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	data := mapper.EncodeGetUserByIdResponse(user)

	c.JSON(200, constants.JsonResp{
		Data:    data,
		Message: "User details fetched successfully",
	})
}
func (ctl *userController) AddFriend(c *gin.Context) {
	ctx := c.Request.Context()
	CurUserId, _ := c.Get("UserID")
	FriendUserId := c.Query("user_id")

	req, err := mapper.DecodeAddFriendRequest(CurUserId, FriendUserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AddFriendResponse{
			Message: "Error decoding user_id: " + err.Error(),
		})
		return
	}
	nerr := ctl.userService.AddFriend(ctx, req)
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
func (ctl *userController) GetFriends(c *gin.Context) {
	ctx := c.Request.Context()
	curUserId, _ := c.Get("UserID")
	uid, ok := curUserId.(uint32)
	if !ok {
		c.JSON(http.StatusBadRequest, dto.GetFriendsResponse{
			Friends: []dto.Friend{},
		})
		return
	}
	req := &dto.GetFriendsRequest{UserID: uid}
	friends, nerr := ctl.userService.GetFriends(ctx, req)
	if nerr != nil {
		c.JSON(int(nerr.Code), constants.JsonResp{
			Message: nerr.Message,
		})
		return
	}
	c.JSON(http.StatusOK, dto.GetFriendsResponse{
		Friends: friends,
	})
}
func (ctl *userController) RemoveFriend(c *gin.Context) {
	ctx := c.Request.Context()
	CurUserId, _ := c.Get("UserID")
	FriendUserId := c.Query("user_id")

	req, err := mapper.DecodeAddFriendRequest(CurUserId, FriendUserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AddFriendResponse{
			Message: "Error decoding user_id: " + err.Error(),
		})
		return
	}
	nerr := ctl.userService.RemoveFriend(ctx, req)
	if nerr != nil {
		c.JSON(int(nerr.Code), constants.JsonResp{
			Message: nerr.Message,
		})
		return
	}
	c.JSON(http.StatusOK, constants.JsonResp{
		Message: "friend removed succesfully",
	})

}

func (ctl *userController) ActivateAllUsers(c *gin.Context) {
	ctx := c.Request.Context()
	err := ctl.userService.ActivateAllUsers(ctx)
	if err != nil {
		c.JSON(int(err.Code), constants.JsonResp{
			Message: err.Message,
		})
		return
	}
	c.JSON(http.StatusOK, constants.JsonResp{
		Message: "All users activated successfully",
	})
}
