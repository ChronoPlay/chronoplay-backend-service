package service

import (
	"context"
	"fmt"
	"log"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/model"
	"github.com/ChronoPlay/chronoplay-backend-service/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardService interface {
	AddCard(ctx context.Context, req dto.AddCardRequest) *helpers.CustomError
	GetCard(ctx context.Context, req dto.GetCardRequest) (res *model.Card, err *helpers.CustomError)
}

type cardService struct {
	cardRepo model.CardRepository
	UserRepo model.UserRepository
}

func NewCardService(cardRepo model.CardRepository, userRepo model.UserRepository) CardService {
	return &cardService{
		cardRepo: cardRepo,
		UserRepo: userRepo,
	}
}

func (s *cardService) AddCard(ctx context.Context, req dto.AddCardRequest) *helpers.CustomError {
	if req.UserId == 0 {
		return helpers.BadRequest("User ID is required")
	}
	users, err := s.UserRepo.GetUsers(ctx, model.User{UserId: req.UserId})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return helpers.NotFound("User not found")
	}
	fmt.Println("User:", users[0])
	req.UserType = users[0].UserType

	err = utils.ValidateAddCardRequest(req)
	if err != nil {
		return err
	}

	// uploading image on cloudinary
	imageUrl, err := utils.UploadImageToCloudinary(ctx, req.Image)
	if err != nil {
		return err
	}

	log.Println("Entered here - RegisterUser (userService)")
	session, derr := s.cardRepo.GetCollection().Database().Client().StartSession()
	if derr != nil {
		return helpers.System("Failed to start session: " + derr.Error())
	}
	defer session.EndSession(ctx)
	log.Println("Successfully created session")

	merr := mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
		// You can use sessCtx instead of ctx for transactional operations

		existingCards, err := s.cardRepo.GetCards(sessCtx, model.GetCardsRequest{
			Number: req.CardNumber,
		})
		if err != nil {
			return err
		}
		if len(existingCards) != 0 {
			return helpers.BadRequest("Card exists with given card number")
		}
		// Call repository method with sessCtx
		err = s.cardRepo.AddCard(sessCtx, model.Card{
			Name:        req.CardName,
			Number:      req.CardNumber,
			Description: req.CardDescription,
			Total:       req.TotalCards,
			Creator:     req.UserId,
			ImageUrl:    imageUrl,
		})
		if err != nil {
			return err
		}

		return nil // Will commit if nil
	})

	if merr != nil {
		return helpers.NoType(merr)
	}

	return nil
}

func (s *cardService) GetCard(ctx context.Context, req dto.GetCardRequest) (res *model.Card, err *helpers.CustomError) {
	err = utils.ValidateGetCardRequest(req)
	if err != nil {
		return res, err
	}

	log.Println("Getting Name by Card Number")
	// Call repository method with sessCtx
	card, err := s.cardRepo.GetCardByNumber(ctx, req.CardNumber)

	if err != nil {
		return res, err
	}
	return card, nil
}
