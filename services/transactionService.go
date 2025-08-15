package service

import (
	"context"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/model"
	"github.com/ChronoPlay/chronoplay-backend-service/utils"
)

type TransactionService interface {
	TransferCash(ctx context.Context, req dto.TransferCashRequest) *helpers.CustomError
	TransferCards(ctx context.Context, req dto.TransferCardRequest) *helpers.CustomError
	GiveCards(ctx context.Context, req dto.TransferCardRequest) *helpers.CustomError
}

type transactionService struct {
	cardTransactionRepo model.CardTransactionRepository
	cashTransactionRepo model.CashTransactionRepository
	userRepo            model.UserRepository
	cardRepo            model.CardRepository
}

func NewTransactionService(cardTransactionRepo model.CardTransactionRepository, cashTransactionRepo model.CashTransactionRepository, userRepo model.UserRepository, cardRepo model.CardRepository) TransactionService {
	return &transactionService{
		cardTransactionRepo: cardTransactionRepo,
		cashTransactionRepo: cashTransactionRepo,
		userRepo:            userRepo,
		cardRepo:            cardRepo,
	}
}

func (s *transactionService) TransferCash(ctx context.Context, req dto.TransferCashRequest) (err *helpers.CustomError) {
	err = utils.ValidateTransferCashRequest(req)
	if err != nil {
		return err
	}

	users, err := s.userRepo.GetUsers(ctx, model.User{UserId: req.GivenBy})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return helpers.NotFound("User not found")
	}
	recieverUsers, err := s.userRepo.GetUsers(ctx, model.User{UserId: req.GivenTo})
	if err != nil {
		return err
	}
	if len(recieverUsers) == 0 {
		return helpers.NotFound("User not found")
	}
	err = s.IsCashTransactionPossible(ctx, dto.IsCashTransactionPossibleRequest{
		User:   users[0],
		Amount: req.Amount,
	})
	if err != nil {
		return err
	}

	if !utils.IsValidTransactionStatus(req.Status) {
		req.Status = model.TRANSACTION_STATUS_PENDING
	}

	transaction := model.CashTransaction{
		Amount:    req.Amount,
		GivenBy:   req.GivenBy,
		GivenTo:   req.GivenTo,
		Status:    req.Status,
		CreatedBy: req.GivenBy,
	}

	session, serr := s.cashTransactionRepo.GetCollection().Database().Client().StartSession()
	if serr != nil {
		return helpers.System("Failed to start session: " + serr.Error())
	}
	defer session.EndSession(ctx)
	serr = session.StartTransaction()
	if serr != nil {
		return helpers.System("Failed to start transaction: " + serr.Error())
	}
	defer func() {
		if err != nil {
			if abortErr := session.AbortTransaction(ctx); abortErr != nil {
				err = helpers.System("Failed to abort transaction: " + abortErr.Error())
			}
		} else {
			if commitErr := session.CommitTransaction(ctx); commitErr != nil {
				err = helpers.System("Failed to commit transaction: " + commitErr.Error())
			}
		}
	}()
	_, err = s.cashTransactionRepo.AddCashTransaction(ctx, transaction)
	if err != nil {
		return err
	}
	if transaction.Status == model.TRANSACTION_STATUS_SUCCESS {
		user := users[0]
		user.Cash = user.Cash - req.Amount
		err = s.userRepo.UpdateUser(ctx, user)
		if err != nil {
			return err
		}
		user = recieverUsers[0]
		user.Cash = user.Cash + req.Amount
		err = s.userRepo.UpdateUser(ctx, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *transactionService) TransferCards(ctx context.Context, req dto.TransferCardRequest) (err *helpers.CustomError) {
	users, err := s.userRepo.GetUsers(ctx, model.User{UserId: req.UserId})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return helpers.NotFound("User not found")
	}
	req.UserType = users[0].UserType
	err = utils.ValidateTransferCardsRequest(req)
	if err != nil {
		return err
	}
	if req.GivenBy != 0 {
		users, err = s.userRepo.GetUsers(ctx, model.User{UserId: req.GivenBy})
		if err != nil {
			return err
		}
		if len(users) == 0 {
			return helpers.NotFound("User not found")
		}
	}
	recieverUsers, err := s.userRepo.GetUsers(ctx, model.User{UserId: req.GivenTo})
	if err != nil {
		return err
	}
	if len(recieverUsers) == 0 {
		return helpers.NotFound("User not found")
	}
	cardNumbers := []string{}
	for _, card := range req.Cards {
		cardNumbers = append(cardNumbers, card.CardNumber)
	}
	cards, err := s.cardRepo.GetCards(ctx, model.GetCardsRequest{
		Numbers: cardNumbers,
	})
	if err != nil {
		return err
	}
	if len(cards) == 0 || len(cards) != len(req.Cards) {
		return helpers.NotFound("Some cards not found")
	}
	cardsToTransferMap := make(map[string]uint32)
	for _, card := range req.Cards {
		cardsToTransferMap[card.CardNumber] = card.Amount
	}
	cardsOccupiedMap := make(map[string]uint32)
	for _, card := range users[0].Cards {
		cardsOccupiedMap[card.CardNumber] = card.Occupied
	}
	cardsMap := make(map[string]model.Card)
	for _, card := range cards {
		cardsMap[card.Number] = card
	}
	err = s.IsCardTransactionPossible(ctx, dto.IsCardTransactionPossibleRequest{
		GivenBy:            req.GivenBy,
		User:               users[0],
		CardsToTransferMap: cardsToTransferMap,
		CardsOccupiedMap:   cardsOccupiedMap,
		CardsMap:           cardsMap,
	})
	if err != nil {
		return err
	}
	if !utils.IsValidTransactionStatus(req.Status) {
		req.Status = model.TRANSACTION_STATUS_PENDING
	}
	transactions := []model.CardTransaction{}
	for _, card := range req.Cards {
		transaction := model.CardTransaction{
			CardNumber: card.CardNumber,
			Amount:     card.Amount,
			GivenBy:    req.GivenBy,
			GivenTo:    req.GivenTo,
			Status:     req.Status,
			CreatedBy:  req.GivenBy,
		}
		transactions = append(transactions, transaction)
	}
	session, serr := s.cardTransactionRepo.GetCollection().Database().Client().StartSession()
	if serr != nil {
		return helpers.System("Failed to start session: " + serr.Error())
	}
	defer session.EndSession(ctx)
	serr = session.StartTransaction()
	if serr != nil {
		return helpers.System("Failed to start transaction: " + serr.Error())
	}
	defer func() {
		if err != nil {
			if abortErr := session.AbortTransaction(ctx); abortErr != nil {
				err = helpers.System("Failed to abort transaction: " + abortErr.Error())
			}
		} else {
			if commitErr := session.CommitTransaction(ctx); commitErr != nil {
				err = helpers.System("Failed to commit transaction: " + commitErr.Error())
			}
		}
	}()
	_, err = s.cardTransactionRepo.AddCardTransactions(ctx, transactions)
	if err != nil {
		return err
	}
	if req.Status == model.TRANSACTION_STATUS_SUCCESS {
		if req.GivenBy != 0 {
			user := users[0]
			for _, card := range req.Cards {
				cardFound := false
				for i, userCard := range user.Cards {
					if userCard.CardNumber == card.CardNumber {
						cardFound = true
						if userCard.Occupied < card.Amount {
							return helpers.BadRequest("Insufficient card balance")
						}
						user.Cards[i].Occupied -= card.Amount
						break
					}
				}
				if !cardFound {
					return helpers.BadRequest("User does not have the required card")
				}
			}
			err = s.userRepo.UpdateUser(ctx, user)
			if err != nil {
				return err
			}
		} else {
			for _, cardToTransfer := range req.Cards {
				cardFound := false
				for i, card := range cards {
					if cardToTransfer.CardNumber == card.Number {
						cardFound = true
						if card.Total-card.Occupied < cardToTransfer.Amount {
							return helpers.BadRequest("Insufficient card balance")
						}
						cards[i].Occupied += cardToTransfer.Amount
						break
					}
				}
				if !cardFound {
					return helpers.BadRequest("card does not exist: " + cardToTransfer.CardNumber)
				}
			}
			err = s.cardRepo.UpdateCards(ctx, cards)
			if err != nil {
				return err
			}
		}
		recieverUser := recieverUsers[0]
		for _, card := range req.Cards {
			cardFound := false
			for i, recieversCard := range recieverUser.Cards {
				if card.CardNumber == recieversCard.CardNumber {
					cardFound = true
					recieverUser.Cards[i].Occupied += card.Amount
					break
				}
			}
			if !cardFound {
				recieverUser.Cards = append(recieverUser.Cards, model.CardOccupied{
					CardNumber: card.CardNumber,
					Occupied:   card.Amount,
				})
			}
		}
		err = s.userRepo.UpdateUser(ctx, recieverUser)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *transactionService) GiveCards(ctx context.Context, req dto.TransferCardRequest) (err *helpers.CustomError) {
	err = utils.ValidateGiveCardsRequest(req)
	if err != nil {
		return err
	}
	users, err := s.userRepo.GetUsers(ctx, model.User{UserId: req.UserId})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return helpers.NotFound("User not found")
	}
	if !utils.IsAdmin(users[0].UserType) {
		return helpers.Unauthorized("Only admin can give cards to users")
	}

	// Implement the logic to give cards
	return nil
}

func (s *transactionService) Exchange(ctx context.Context, req dto.ExchangeRequest) {
}

func (s *transactionService) GetTransactions(ctx context.Context, req dto.GetTransactionsRequest) {
}

func (s *transactionService) IsCashTransactionPossible(ctx context.Context, req dto.IsCashTransactionPossibleRequest) *helpers.CustomError {
	// for now here is only one check but later more checks may be added
	if req.User.Cash < req.Amount {
		return helpers.BadRequest("Insufficient cash")
	}
	return nil
}

func (s *transactionService) IsCardTransactionPossible(ctx context.Context, req dto.IsCardTransactionPossibleRequest) *helpers.CustomError {
	if req.GivenBy != 0 {
		for cardNumber, amount := range req.CardsToTransferMap {
			if occupied, exists := req.CardsOccupiedMap[cardNumber]; exists {
				if occupied < amount {
					return helpers.BadRequest("Insufficient card balance for card: " + cardNumber)
				}
			} else {
				return helpers.BadRequest("Card not found: " + cardNumber)
			}
		}
	} else {
		for cardNumber, amount := range req.CardsToTransferMap {
			if card, exists := req.CardsMap[cardNumber]; exists {
				if card.Total-card.Occupied < amount {
					return helpers.BadRequest("Insufficient card balance for card: " + cardNumber)
				}
			} else {
				return helpers.BadRequest("Card not found: " + cardNumber)
			}
		}
	}
	return nil
}
