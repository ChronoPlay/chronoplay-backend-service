package service

import (
	"context"
	"sort"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	"github.com/ChronoPlay/chronoplay-backend-service/mapper"
	model "github.com/ChronoPlay/chronoplay-backend-service/model"
	"github.com/ChronoPlay/chronoplay-backend-service/utils"
)

type TransactionService interface {
	TransferCash(ctx context.Context, req dto.TransferCashRequest) *helpers.CustomError
	TransferCards(ctx context.Context, req dto.TransferCardRequest) *helpers.CustomError
	GiveCards(ctx context.Context, req dto.TransferCardRequest) *helpers.CustomError
	GetTransactions(ctx context.Context, req dto.GetTransactionsRequest) (dto.GetTransactionsResponse, *helpers.CustomError)
	Exchange(ctx context.Context, req dto.ExchangeRequest) *helpers.CustomError
	GetPossibleExchange(ctx context.Context, req dto.GetPossibleExchangeRequest) (dto.GetPossibleExchangeResponse, *helpers.CustomError)
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
	users, err := s.userRepo.GetUsers(ctx, model.User{UserId: req.UserId})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return helpers.NotFound("User not found")
	}
	req.UserType = users[0].UserType
	err = utils.ValidateTransferCashRequest(req)
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
	err = s.IsCashTransactionPossible(ctx, dto.IsCashTransactionPossibleRequest{
		GivenBy: req.GivenBy,
		User:    users[0],
		Amount:  req.Amount,
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
		var user model.User
		if req.GivenBy != 0 {
			user = users[0]
			user.Cash = user.Cash - req.Amount
			err = s.userRepo.UpdateUser(ctx, user)
			if err != nil {
				return err
			}
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

func (s *transactionService) Exchange(ctx context.Context, req dto.ExchangeRequest) *helpers.CustomError {
	err := utils.ValidateExchangeRequest(req)
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
	sender := users[0]
	users, err = s.userRepo.GetUsers(ctx, model.User{UserId: req.GivenTo})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return helpers.NotFound("User not found")
	}
	receiver := users[0]
	if sender.UserId == receiver.UserId {
		return helpers.BadRequest("Sender and receiver cannot be the same user")
	}

	err = IsExchangePossible(dto.IsExhangePossibleRequest{
		GivenByUser:   sender,
		GivenToUser:   receiver,
		CashSent:      req.CashSent,
		CashRecieved:  req.CashRecieved,
		CardsSent:     req.CardsSent,
		CardsRecieved: req.CardsRecieved,
	})
	if err != nil {
		return err
	}

	cardTransactions := []model.CardTransaction{}
	for _, card := range req.CardsSent {
		cardTransactions = append(cardTransactions, model.CardTransaction{
			CardNumber: card.CardNumber,
			Amount:     card.Amount,
			GivenBy:    sender.UserId,
			GivenTo:    receiver.UserId,
			Status:     model.TRANSACTION_STATUS_PENDING,
			CreatedBy:  sender.UserId,
		})
	}
	for _, card := range req.CardsRecieved {
		cardTransactions = append(cardTransactions, model.CardTransaction{
			CardNumber: card.CardNumber,
			Amount:     card.Amount,
			GivenBy:    receiver.UserId,
			GivenTo:    sender.UserId,
			Status:     model.TRANSACTION_STATUS_PENDING,
			CreatedBy:  receiver.UserId,
		})
	}

	// need to add transaction commit and abort and rollback handling here
	paymentGuid := uint32(0)
	if len(cardTransactions) > 0 {
		paymentGuid, err = s.cardTransactionRepo.AddCardTransactions(ctx, cardTransactions)
		if err != nil {
			return err
		}
	}

	cashTransaction := model.CashTransaction{}
	if req.CashSent > 0 {
		cashTransaction = model.CashTransaction{
			Amount:          req.CashSent,
			GivenBy:         sender.UserId,
			GivenTo:         receiver.UserId,
			Status:          model.TRANSACTION_STATUS_PENDING,
			CreatedBy:       sender.UserId,
			TransactionGuid: paymentGuid,
		}
	} else if req.CashRecieved > 0 {
		cashTransaction = model.CashTransaction{
			Amount:          req.CashRecieved,
			GivenBy:         receiver.UserId,
			GivenTo:         sender.UserId,
			Status:          model.TRANSACTION_STATUS_PENDING,
			CreatedBy:       receiver.UserId,
			TransactionGuid: paymentGuid,
		}
	}
	_, err = s.cashTransactionRepo.AddCashTransaction(ctx, cashTransaction)
	if err != nil {
		return err
	}
	return nil
}

func (s *transactionService) GetTransactions(ctx context.Context, req dto.GetTransactionsRequest) (resp dto.GetTransactionsResponse, err *helpers.CustomError) {
	if req.UserId == 0 {
		return resp, helpers.BadRequest("User ID is required")
	}
	// first i need to get all cashtransactions which are given by recieved by this user
	cashTransactionsByUser, err := s.cashTransactionRepo.GetCashTransactionsByUserId(ctx, req.UserId)
	if err != nil {
		return resp, err
	}
	cashTransactionsToUser, err := s.cashTransactionRepo.GetCashTransactionsToUserId(ctx, req.UserId)
	if err != nil {
		return resp, err
	}

	// then i need to map those by transaction guid
	guidToCashTransactionsByUserMap := make(map[uint32][]model.CashTransaction)
	for _, transaction := range cashTransactionsByUser {
		guidToCashTransactionsByUserMap[transaction.TransactionGuid] = append(guidToCashTransactionsByUserMap[transaction.TransactionGuid], transaction)
	}
	guidToCashTransactionsToUserMap := make(map[uint32][]model.CashTransaction)
	for _, transaction := range cashTransactionsToUser {
		guidToCashTransactionsToUserMap[transaction.TransactionGuid] = append(guidToCashTransactionsToUserMap[transaction.TransactionGuid], transaction)
	}

	// then i need to get all card transactions which are given by recieved by this user
	cardTransactionsByUser, err := s.cardTransactionRepo.GetCardTransactionsByUserId(ctx, req.UserId)
	if err != nil {
		return resp, err
	}
	cardTransactionsToUser, err := s.cardTransactionRepo.GetCardTransactionsToUserId(ctx, req.UserId)
	if err != nil {
		return resp, err
	}

	// then i need to map those by transaction guid
	guidToCardTransactionsByUserMap := make(map[uint32][]model.CardTransaction)
	for _, transaction := range cardTransactionsByUser {
		guidToCardTransactionsByUserMap[transaction.TransactionGuid] = append(guidToCardTransactionsByUserMap[transaction.TransactionGuid], transaction)
	}
	guidToCardTransactionsToUserMap := make(map[uint32][]model.CardTransaction)
	for _, transaction := range cardTransactionsToUser {
		guidToCardTransactionsToUserMap[transaction.TransactionGuid] = append(guidToCardTransactionsToUserMap[transaction.TransactionGuid], transaction)
	}

	// after this i need to merge those maps on guid basis and return my response
	uniqueGuids := make(map[uint32]bool)
	for guid := range guidToCashTransactionsByUserMap {
		uniqueGuids[guid] = true
	}
	for guid := range guidToCardTransactionsByUserMap {
		uniqueGuids[guid] = true
	}
	for guid := range guidToCashTransactionsToUserMap {
		uniqueGuids[guid] = true
	}
	for guid := range guidToCardTransactionsToUserMap {
		uniqueGuids[guid] = true
	}
	for guid := range uniqueGuids {
		transaction := dto.Transaction{
			TransactionGuid: guid,
		}
		cardTransactionsToUser, ok := guidToCardTransactionsToUserMap[guid]
		if ok {
			for _, cardTransaction := range cardTransactionsToUser {
				transaction.CardsRecieved = append(transaction.CardsRecieved, dto.Card{
					CardNumber: cardTransaction.CardNumber,
					Amount:     cardTransaction.Amount,
				})
			}
			transaction.Time = cardTransactionsToUser[0].CreatedAt.Time()
			transaction.TransactionWith = cardTransactionsToUser[0].GivenBy
			transaction.Status = cardTransactionsToUser[0].Status
		}
		cardTransactionsByUser, ok := guidToCardTransactionsByUserMap[guid]
		if ok {
			for _, cardTransaction := range cardTransactionsByUser {
				transaction.CardsSent = append(transaction.CardsSent, dto.Card{
					CardNumber: cardTransaction.CardNumber,
					Amount:     cardTransaction.Amount,
				})
			}
			if transaction.Time.IsZero() {
				transaction.Time = cardTransactionsByUser[0].CreatedAt.Time()
			}
			if transaction.TransactionWith == 0 {
				transaction.TransactionWith = cardTransactionsByUser[0].GivenTo
			}
			if transaction.Status == "" {
				transaction.Status = cardTransactionsByUser[0].Status
			}
		}
		cashTransactionsByUser, ok := guidToCashTransactionsByUserMap[guid]
		if ok {
			for _, cashTransaction := range cashTransactionsByUser {
				transaction.CashSent += cashTransaction.Amount
			}
			if transaction.Time.IsZero() {
				transaction.Time = cashTransactionsByUser[0].CreatedAt.Time()
			}
			if transaction.TransactionWith == 0 {
				transaction.TransactionWith = cashTransactionsByUser[0].GivenTo
			}
			if transaction.Status == "" {
				transaction.Status = cashTransactionsByUser[0].Status
			}
		}
		cashTransactionsToUser, ok := guidToCashTransactionsToUserMap[guid]
		if ok {
			for _, cashTransaction := range cashTransactionsToUser {
				transaction.CashRecieved += cashTransaction.Amount
			}
			if transaction.Time.IsZero() {
				transaction.Time = cashTransactionsToUser[0].CreatedAt.Time()
			}
			if transaction.TransactionWith == 0 {
				transaction.TransactionWith = cashTransactionsToUser[0].GivenBy
			}
			if transaction.Status == "" {
				transaction.Status = cashTransactionsToUser[0].Status
			}
		}
		resp.Transactions = append(resp.Transactions, transaction)
	}
	sort.Slice(resp.Transactions, func(i, j int) bool {
		return resp.Transactions[i].TransactionGuid > resp.Transactions[j].TransactionGuid
	})
	return resp, nil
}

func (s *transactionService) IsCashTransactionPossible(ctx context.Context, req dto.IsCashTransactionPossibleRequest) *helpers.CustomError {
	// for now here is only one check but later more checks may be added
	if req.GivenBy == 0 {
		// amount is being paid by system
		return nil
	}
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

func IsExchangePossible(req dto.IsExhangePossibleRequest) *helpers.CustomError {
	err := IsValidCardExchange(req.CardsSent, req.CardsRecieved)
	if err != nil {
		return err
	}

	if len(req.CardsSent) > 0 {
		cardsToSendMap := make(map[string]uint32)
		for _, card := range req.CardsSent {
			if card.Amount <= 0 {
				return helpers.BadRequest("Card amount must be greater than zero")
			}
			cardsToSendMap[card.CardNumber] += card.Amount
		}
		cardsSenderHaveMap := make(map[string]uint32)
		for _, card := range req.GivenByUser.Cards {
			cardsSenderHaveMap[card.CardNumber] += card.Occupied
		}
		for cardNumber, amount := range cardsToSendMap {
			if cardsSenderHaveMap[cardNumber] < amount {
				return helpers.BadRequest("Insufficient balance for card: " + cardNumber)
			}
		}
	}
	if len(req.CardsRecieved) > 0 {
		cardsToReceiveMap := make(map[string]uint32)
		for _, card := range req.CardsRecieved {
			if card.Amount <= 0 {
				return helpers.BadRequest("Card amount must be greater than zero")
			}
			cardsToReceiveMap[card.CardNumber] += card.Amount
		}
		cardsReceiverHaveMap := make(map[string]uint32)
		for _, card := range req.GivenToUser.Cards {
			cardsReceiverHaveMap[card.CardNumber] += card.Occupied
		}
		for cardNumber, amount := range cardsToReceiveMap {
			if cardsReceiverHaveMap[cardNumber] < amount {
				return helpers.BadRequest("Insufficient balance for card: " + cardNumber)
			}
		}
	}
	if req.CashSent > req.GivenByUser.Cash {
		return helpers.BadRequest("Insufficient cash to exchange")
	}
	if req.CashRecieved > req.GivenToUser.Cash {
		return helpers.BadRequest("Insufficient cash to exchange")
	}
	return nil
}

func IsValidCardExchange(cardsSent, cardsRecieved []dto.Card) *helpers.CustomError {
	sentCardNumbers := make(map[string]bool)
	for _, card := range cardsSent {
		sentCardNumbers[card.CardNumber] = true
	}

	for _, card := range cardsRecieved {
		if sentCardNumbers[card.CardNumber] {
			return helpers.BadRequest("Same card cannot be sent and received")
		}
	}
	return nil
}

func (s *transactionService) GetPossibleExchange(ctx context.Context, req dto.GetPossibleExchangeRequest) (dto.GetPossibleExchangeResponse, *helpers.CustomError) {
	err := utils.ValidateGetPossibleExchangeRequest(req)
	if err != nil {
		return dto.GetPossibleExchangeResponse{}, err
	}
	users, err := s.userRepo.GetUsers(ctx, model.User{UserId: req.UserId})
	if err != nil {
		return dto.GetPossibleExchangeResponse{}, err
	}
	if len(users) == 0 {
		return dto.GetPossibleExchangeResponse{}, helpers.NotFound("User not found")
	}
	curUser := users[0]
	traderUsers, err := s.userRepo.GetUsers(ctx, model.User{UserId: req.TraderId})
	if err != nil {
		return dto.GetPossibleExchangeResponse{}, err
	}
	if len(traderUsers) == 0 {
		return dto.GetPossibleExchangeResponse{}, helpers.NotFound("Trader not found")
	}
	trader := traderUsers[0]

	curUserCardNumbers := []string{}
	for _, card := range curUser.Cards {
		curUserCardNumbers = append(curUserCardNumbers, card.CardNumber)
	}
	traderCardNumbers := []string{}
	for _, card := range trader.Cards {
		traderCardNumbers = append(traderCardNumbers, card.CardNumber)
	}
	curUserCards, err := s.cardRepo.GetCards(ctx, model.GetCardsRequest{
		Numbers: curUserCardNumbers,
	})
	if err != nil {
		return dto.GetPossibleExchangeResponse{}, err
	}
	traderCards, err := s.cardRepo.GetCards(ctx, model.GetCardsRequest{
		Numbers: traderCardNumbers,
	})
	if err != nil {
		return dto.GetPossibleExchangeResponse{}, err
	}

	return dto.GetPossibleExchangeResponse{
		YourCash:    curUser.Cash,
		TraderCash:  trader.Cash,
		YourCards:   mapper.MapCardsToResponse(curUserCards),
		TraderCards: mapper.MapCardsToResponse(traderCards),
	}, nil
}
