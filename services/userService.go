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
	GetUser(context.Context, model.User) (*model.User, *helpers.CustomError)
	RegisterUser(ctx context.Context, req model.User) (err *helpers.CustomError)
	VerifyUser(ctx context.Context, req dto.VerifyUserRequest) (err *helpers.CustomError)
	LoginUser(ctx context.Context, req dto.LoginUserRequest) (dto.LoginUserResponse, *helpers.CustomError)
}

type userService struct {
	userRepo model.UserRepository
}

func NewUserService(userRepo model.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUser(ctx context.Context, req model.User) (resp *model.User, err *helpers.CustomError) {
	users, err := s.userRepo.GetUsers(ctx, req)
	if err != nil {
		return resp, err
	}
	if len(users) == 0 {
		return nil, helpers.NotFound("User not found")
	}
	return &users[0], nil
}

func (s *userService) RegisterUser(ctx context.Context, req model.User) (err *helpers.CustomError) {
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
		})
		if err != nil {
			return err
		}
		if len(existingUsers) != 0 {
			return helpers.BadRequest("User exists with given userName")
		}
		existingUsers, err = s.userRepo.GetUsers(sessCtx, model.User{
			Email: req.Email,
		})
		if err != nil {
			return err
		}
		if len(existingUsers) != 0 {
			return helpers.BadRequest("User exists with given emailId")
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
		return helpers.NoType(merr)
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

func (s *userService) VerifyUser(ctx context.Context, req dto.VerifyUserRequest) (err *helpers.CustomError) {
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

func (s *userService) LoginUser(ctx context.Context, req dto.LoginUserRequest) (resp dto.LoginUserResponse, err *helpers.CustomError) {
	log.Println("LoginUser: Starting login process with req body: ", req)

	if req.Email == "" && req.UserName == "" {
		log.Println("LoginUser: Missing email or username")
		return resp, helpers.BadRequest("Email or username is required")
	}
	if req.Password == "" {
		log.Println("LoginUser: Missing password")
		return resp, helpers.BadRequest("Password is required")
	}

	log.Println("LoginUser: Fetching user from repository")
	users, err := s.userRepo.GetUsers(ctx, model.User{
		Email:    req.Email,
		UserName: req.UserName,
	})

	if err != nil {
		log.Println("LoginUser: Error fetching user from repository:", err)
		return resp, helpers.System(err.Error())
	}
	if len(users) == 0 {
		log.Println("LoginUser: No user found with given credentials")
		return resp, helpers.BadRequest("User not found with given credentials")
	}

	if !users[0].IsAuthorized {
		log.Println("LoginUser: User is not verified")
		return resp, helpers.Unauthorized("User is not verified yet. Please verify your emailId first")
	}

	log.Println("LoginUser: Verifying password")
	err = utils.CheckPasswordHash(req.Password, users[0].Password)
	if err != nil {
		log.Println("LoginUser: Invalid password")
		return resp, helpers.Unauthorized("Invalid password")
	}

	log.Println("LoginUser: Generating JWT token for user:", users[0].UserId)
	jwtToken, err := utils.GenerateJwtToken(users[0].UserId)
	if err != nil {
		log.Println("LoginUser: Failed to generate JWT token:", err)
		return resp, helpers.System("Failed to generate JWT token: " + err.Error())
	}

	log.Println("LoginUser: User logged in successfully:", users[0].UserId)
	return dto.LoginUserResponse{
		Token: jwtToken,
	}, nil
}
