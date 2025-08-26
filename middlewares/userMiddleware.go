package middleware

import (
	"context"
	"log"
	"time"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/model"
	service "github.com/ChronoPlay/chronoplay-backend-service/services"
)

type userMiddleware struct {
	next service.UserService
}

type UserMiddleware func(service.UserService) service.UserService

func NewUserMiddleware(next service.UserService) UserMiddleware {
	return func(next service.UserService) service.UserService {
		return &userMiddleware{
			next: next,
		}
	}
}

func (mw userMiddleware) GetUser(ctx context.Context, req dto.GetUserRequest) (resp dto.GetUserResponse, err *helpers.CustomError) {
	defer func(begin time.Time) {
		log.Printf("ctx:%v method:%v req:%v took:%v err:%v resp:%v",
			ctx, "GetUser", req, time.Since(begin), err, resp)
	}(time.Now())
	return mw.next.GetUser(ctx, req)
}

func (mw userMiddleware) RegisterUser(ctx context.Context, req model.User) (err *helpers.CustomError) {
	defer func(begin time.Time) {
		log.Printf("ctx:%v method:%v req:%v took:%v err:%v",
			ctx, "RegisterUser", req, time.Since(begin), err)
	}(time.Now())
	return mw.next.RegisterUser(ctx, req)
}

func (mw userMiddleware) VerifyUser(ctx context.Context, req dto.VerifyUserRequest) (err *helpers.CustomError) {
	defer func(begin time.Time) {
		log.Printf("ctx:%v method:%v req:%v took:%v err:%v",
			ctx, "VerifyUser", req, time.Since(begin), err)
	}(time.Now())
	return mw.next.VerifyUser(ctx, req)
}

func (mw userMiddleware) LoginUser(ctx context.Context, req dto.LoginUserRequest) (resp dto.LoginUserResponse, err *helpers.CustomError) {
	defer func(begin time.Time) {
		log.Printf("ctx:%v method:%v req:%v took:%v err:%v resp:%v",
			ctx, "LoginUser", req, time.Since(begin), err, resp)
	}(time.Now())
	return mw.next.LoginUser(ctx, req)
}
func (mw userMiddleware) AddFriend(ctx context.Context, req *dto.AddFriendRequest) (err *helpers.CustomError) {
	defer func(begin time.Time) {
		log.Printf("ctx:%v method:%v userID:%d friendID:%d took:%v err:%v",
			ctx, "AddFriend", req.UserID, req.FriendID, time.Since(begin), err)
	}(time.Now())

	return mw.next.AddFriend(ctx, req)
}

func (mw userMiddleware) GetFriends(ctx context.Context, req *dto.GetFriendsRequest) (resp []dto.Friend, err *helpers.CustomError) {
	defer func(begin time.Time) {
		log.Printf("ctx:%v method:%v userID:%d took:%v err:%v",
			ctx, "GetFriends", req.UserID, time.Since(begin), err)
	}(time.Now())

	return mw.next.GetFriends(ctx, req)
}
func (mw userMiddleware) RemoveFriend(ctx context.Context, req *dto.AddFriendRequest) (err *helpers.CustomError) {
	defer func(begin time.Time) {
		log.Printf("ctx: %v method: AddFriend userID: %d fid: %d took: %s err: %v",
			ctx, req.UserID, req.FriendID, time.Since(begin), err)
	}(time.Now())

	return mw.next.RemoveFriend(ctx, req)
}

func (mw userMiddleware) GetAllActiveUsers() (users []model.User, err *helpers.CustomError) {
	defer func(begin time.Time) {
		log.Printf(" method: GetAllActiveUsers took: %s err: %v",
			time.Since(begin), err)
	}(time.Now())

	return mw.next.GetAllActiveUsers()
}

func (mw userMiddleware) UpdateUser(ctx context.Context, user model.User) (err *helpers.CustomError) {
	defer func(begin time.Time) {
		log.Printf("ctx:%v method:%v req:%v took:%v err:%v",
			ctx, "UpdateUser", user, time.Since(begin), err)
	}(time.Now())
	return mw.next.UpdateUser(ctx, user)
}

func (mw userMiddleware) ActivateAllUsers(ctx context.Context) (err *helpers.CustomError) {
	defer func(begin time.Time) {
		log.Printf("ctx:%v method:%v took:%v err:%v",
			ctx, "ActivateAllUsers", time.Since(begin), err)
	}(time.Now())
	return mw.next.ActivateAllUsers(ctx)
}
