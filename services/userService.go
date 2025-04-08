package service

import (
	"context"

	"github.com/ChronoPlay/chronoplay-backend-service/database"
	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/models"
	"github.com/ChronoPlay/chronoplay-backend-service/utils"
)

type UserService interface {
	GetUser(context.Context, model.User) (model.User, *helpers.CustomEror)
	RegisterUser(ctx context.Context, req model.User) (err *helpers.CustomEror)
}

type userService struct {
	userRepo model.UserRepository
}

func NewUserService(userRepo model.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUser(ctx context.Context, req model.User) (resp model.User, err *helpers.CustomEror) {
	if len(req.UserName) != 0 {
		resp, err = s.userRepo.FindByUserName(req.UserName)
	}
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *userService) RegisterUser(ctx context.Context, req model.User) (err *helpers.CustomEror) {
	err = utils.ValidateUser(req)
	if err != nil {
		return err
	}

	tx := database.DB.Begin()

	req.Password, err = utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = s.userRepo.RegisterUser(tx, req)
	if err != nil {
		return err
	}

	emailVerificationLink := ""
	err = utils.SendEmailToUser(dto.EmailVerificationRequest{
		Email:    req.Email,
		UserName: req.UserName,
		Link:     emailVerificationLink,
	})
	return nil
}
