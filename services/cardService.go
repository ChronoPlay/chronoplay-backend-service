package service

import (
	"context"
	"log"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/model"
	"github.com/ChronoPlay/chronoplay-backend-service/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardService interface {
	AddCard(ctx context.Context, req dto.AddCardRequest) *helpers.CustomError
}

type cardService struct {
	cardRepo model.CardRepository
}

func NewCardService(cardRepo model.CardRepository) CardService {
	return &cardService{
		cardRepo: cardRepo,
	}
}

func (s *cardService) AddCard(ctx context.Context, req dto.AddCardRequest) *helpers.CustomError {
	err := utils.ValidateAddCardRequest(req)
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

		existingCards, err := s.cardRepo.GetCards(sessCtx, model.Card{
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
			Number:      req.CardNumber,
			Description: req.CardDescription,
			Available:   req.TotalCards,
			Creator:     req.UserId,
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
