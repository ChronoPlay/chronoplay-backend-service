package service

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/model"
	"github.com/ChronoPlay/chronoplay-backend-service/utils"
)

type UserService interface {
	GetUser(context.Context, model.User) (*model.User, *helpers.CustomEror)
	RegisterUser(ctx context.Context, req model.User) (err *helpers.CustomEror)
	VerifyUser(ctx context.Context, req model.VerifyUserRequest) (err *helpers.CustomEror)
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

	// user will be unverified in begining
	req.IsAuthorized = false

	log.Println("Entered here - RegisterUser (userService)")
	session, derr := s.userRepo.GetCollection().Database().Client().StartSession()
	if derr != nil {
		return helpers.System("Failed to start session: " + derr.Error())
	}
	defer session.EndSession(ctx)
	log.Println("Successfully created session")

	merr := mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
		// You can use sessCtx instead of ctx for transactional operations

		existingUsers, err := s.userRepo.GetUsers(sessCtx, model.User{
			UserName: req.UserName,
			Email:    req.Email,
		})
		if err != nil {
			return err
		}
		if len(existingUsers) != 0 {
			return helpers.BadRequest("User exists with given userName or emailId")
		}

		req.Password, err = utils.HashPassword(req.Password)
		if err != nil {
			return err // Will abort the transaction
		}

		// Call repository method with sessCtx
		_, err = s.userRepo.RegisterUser(sessCtx, req)
		if err != nil {
			return err
		}

		return nil // Will commit if nil
	})

	if merr != nil {
		return helpers.System("Transaction failed: " + merr.Error())
	}

	emailVerificationLink := utils.GenrateEmailVerificationLink(req.Email)
	err = utils.SendEmailToUser(dto.EmailVerificationRequest{
		Email:    req.Email,
		UserName: req.UserName,
		Link:     emailVerificationLink,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) VerifyUser(ctx context.Context, req model.VerifyUserRequest) (err *helpers.CustomEror) {
	if req.Email == "" {
		return helpers.BadRequest("Email is required")
	}
	user := model.User{
		Email: req.Email,
	}
	existingUsers, err := s.userRepo.GetUsers(ctx, user)
	if err != nil {
		return helpers.System(err.Error())
	}
	if len(existingUsers) == 0 {
		return helpers.BadRequest("No existing user present for this email")
	}

	existingUser := existingUsers[0]
	existingUser.IsAuthorized = true

	err = s.userRepo.UpdateUser(ctx, existingUser)
	if err != nil {
		return helpers.System(err.Error())
	}
	return nil
}
