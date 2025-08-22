package service

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	"github.com/ChronoPlay/chronoplay-backend-service/mapper"
	model "github.com/ChronoPlay/chronoplay-backend-service/model"
	"github.com/ChronoPlay/chronoplay-backend-service/utils"
)

type UserService interface {
	GetUser(context.Context, dto.GetUserRequest) (dto.GetUserResponse, *helpers.CustomError)
	RegisterUser(ctx context.Context, req model.User) (err *helpers.CustomError)
	VerifyUser(ctx context.Context, req dto.VerifyUserRequest) (err *helpers.CustomError)
	LoginUser(ctx context.Context, req dto.LoginUserRequest) (dto.LoginUserResponse, *helpers.CustomError)
	AddFriend(ctx context.Context, req *dto.AddFriendRequest) *helpers.CustomError
	GetFriends(ctx context.Context, req *dto.GetFriendsRequest) ([]dto.Friend, *helpers.CustomError)
	RemoveFriend(ctx context.Context, req *dto.AddFriendRequest) *helpers.CustomError
}

type userService struct {
	userRepo model.UserRepository
	cardRepo model.CardRepository
}

func NewUserService(userRepo model.UserRepository, cardRepo model.CardRepository) UserService {
	return &userService{
		userRepo: userRepo,
		cardRepo: cardRepo,
	}
}

func (s *userService) GetUser(ctx context.Context, req dto.GetUserRequest) (resp dto.GetUserResponse, err *helpers.CustomError) {
	users, err := s.userRepo.GetUsers(ctx, model.User{
		UserId: req.UserID,
	})
	if err != nil {
		return resp, err
	}
	if len(users) == 0 {
		return dto.GetUserResponse{}, helpers.NotFound("User not found")
	}
	cardNumbers := []string{}
	for _, card := range users[0].Cards {
		cardNumbers = append(cardNumbers, card.CardNumber)
	}
	cards, err := s.cardRepo.GetCards(ctx, model.GetCardsRequest{Numbers: cardNumbers})
	if err != nil {
		return dto.GetUserResponse{}, err
	}

	return dto.GetUserResponse{
		Name:        users[0].Name,
		Email:       users[0].Email,
		UserName:    users[0].UserName,
		Cash:        users[0].Cash,
		FriendIds:   users[0].Friends,
		PhoneNumber: users[0].PhoneNumber,
		Cards:       mapper.MapCardsToResponse(cards),
		UserType:    users[0].UserType,
	}, nil
}

func (s *userService) RegisterUser(ctx context.Context, req model.User) (err *helpers.CustomError) {
	err = utils.ValidateUser(req)
	fmt.Println("req:", req)
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

// add friend
func (s *userService) AddFriend(ctx context.Context, req *dto.AddFriendRequest) *helpers.CustomError {
	if req.UserID == req.FriendID {
		return helpers.BadRequest("user cannot add themselves as a friend")
	}
	curUsers, err := s.userRepo.GetUsers(ctx, model.User{ //curUser now have current user data
		UserId: req.UserID,
	})
	if err != nil {
		return err
	}
	if len(curUsers) == 0 {
		return helpers.NotFound("current user not found")
	}
	curUser := curUsers[0]
	friends, ferr := s.userRepo.GetUsers(ctx, model.User{ //friend now have to be added user's data
		UserId: req.FriendID,
	})
	if ferr != nil {
		return ferr
	}
	if len(friends) == 0 {
		return helpers.BadRequest("the user doesn't exists")
	}
	curUserFriends := curUser.Friends
	isfound := false
	for _, friend := range curUserFriends {
		if friend == req.FriendID {
			isfound = true
			break
		}
	}
	if isfound {
		return helpers.BadRequest("user is already a friend with current user")
	}
	curUserFriends = append(curUserFriends, req.FriendID)
	curUser.Friends = curUserFriends
	ferr = s.userRepo.UpdateUser(ctx, curUser)
	if ferr != nil {
		return ferr
	}

	return nil
}
func (s *userService) GetFriends(ctx context.Context, req *dto.GetFriendsRequest) ([]dto.Friend, *helpers.CustomError) {
	users, err := s.userRepo.GetUsers(ctx, model.User{
		UserId: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, helpers.NotFound("user not found ")
	}
	curUser := users[0]
	var friends []dto.Friend
	for _, fid := range curUser.Friends {
		friendUsers, ferr := s.userRepo.GetUsers(ctx, model.User{
			UserId: fid,
		})
		if ferr != nil {
			return nil, ferr
		}
		if len(friendUsers) == 0 {
			continue
		}
		friend := friendUsers[0]
		friends = append(friends, dto.Friend{
			UserID:   friend.UserId,
			UserName: friend.UserName,
			Email:    friend.Email,
		})
	}

	return friends, nil
}

func (s *userService) RemoveFriend(ctx context.Context, req *dto.AddFriendRequest) *helpers.CustomError {
	if req.UserID == req.FriendID {
		return helpers.BadRequest("user cannot remove themselves as a friend")
	}
	curUsers, err := s.userRepo.GetUsers(ctx, model.User{ //curUser now have current user data
		UserId: req.UserID,
	})
	if err != nil {
		return err
	}
	if len(curUsers) == 0 {
		return helpers.NotFound("current user not found")
	}
	curUser := curUsers[0]
	friends, ferr := s.userRepo.GetUsers(ctx, model.User{ //friend now have to be added user's data
		UserId: req.FriendID,
	})
	if ferr != nil {
		return ferr
	}
	if len(friends) == 0 {
		return helpers.BadRequest("the user doesn't exists")
	}
	curUserFriends := curUser.Friends
	isfound := false
	newFriends := []uint32{}

	for _, friend := range curUserFriends {
		if friend == req.FriendID {
			isfound = true
			continue
		}
		newFriends = append(newFriends, friend)
	}
	if !isfound {
		return helpers.BadRequest("user is not a friend with current user")
	}
	curUser.Friends = newFriends
	ferr = s.userRepo.UpdateUser(ctx, curUser)
	if ferr != nil {
		return ferr
	}
	return nil
}
