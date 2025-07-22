package service

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/models"
	"github.com/ChronoPlay/chronoplay-backend-service/utils"
)

type UserService interface {
	GetUser(context.Context, model.User) (*model.User, *helpers.CustomEror)
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

func (s *userService) GetUser(ctx context.Context, req model.User) (resp *model.User, err *helpers.CustomEror) {
	if len(req.UserName) != 0 {
		resp, err = s.userRepo.FindByUserName(ctx, req.UserName)
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

	log.Println("Entered here - RegisterUser (userService)")
	session, derr := s.userRepo.GetCollection().Database().Client().StartSession()
	if derr != nil {
		return helpers.System("Failed to start session: " + derr.Error())
	}
	defer session.EndSession(ctx)
	log.Println("Successfully created session")

	merr := mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
		// You can use sessCtx instead of ctx for transactional operations

		req.Password, err = utils.HashPassword(req.Password)
		if err != nil {
			return err // Will abort the transaction
		}

		// Call repository method with sessCtx
		err = s.userRepo.RegisterUser(sessCtx, req)
		if err != nil {
			return err
		}

		emailVerificationLink := "" // Ideally generate a real one
		err = utils.SendEmailToUser(dto.EmailVerificationRequest{
			Email:    req.Email,
			UserName: req.UserName,
			Link:     emailVerificationLink,
		})

		return err // Will commit if nil
	})

	if merr != nil {
		return helpers.System("Transaction failed: " + merr.Error())
	}

	return nil
}
