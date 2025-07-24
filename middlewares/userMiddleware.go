package middleware

import (
	"context"
	"log"
	"time"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/models"
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

func (mw userMiddleware) GetUser(ctx context.Context, req model.User) (resp *model.User, err *helpers.CustomEror) {
	defer func(begin time.Time) {
		log.Printf("ctx:", ctx, " method:", "GetUser", " req:", req, " took:", time.Since(begin), " err:", err, " resp:", resp)
	}(time.Now())
	return mw.next.GetUser(ctx, req)
}

func (mw userMiddleware) RegisterUser(ctx context.Context, req model.User) (err *helpers.CustomEror) {
	defer func(begin time.Time) {
		log.Printf("ctx:", ctx, " method:", "RegisterUser", " req:", req, " took:", time.Since(begin), " err:", err)
	}(time.Now())
	return mw.next.RegisterUser(ctx, req)
}

func (mw userMiddleware) VerifyUser(ctx context.Context, req model.VerifyUserRequest) (err *helpers.CustomEror) {
	defer func(begin time.Time) {
		log.Printf("ctx:", ctx, " method:", "RegisterUser", " req:", req, " took:", time.Since(begin), " err:", err)
	}(time.Now())
	return mw.next.VerifyUser(ctx, req)
}
