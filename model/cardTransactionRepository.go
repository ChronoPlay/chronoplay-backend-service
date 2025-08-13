package model

import (
	"context"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardTransaction struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TransactionId   uint32             `bson:"transaction_id" json:"transaction_id"`
	TransactionGuid uint32             `bson:"transaction_guid" json:"transaction_guid"`
	Amount          float32            `bson:"amount" json:"amount"`
	CardNumber      string             `bson:"card_number" json:"card_number"`
	GivenBy         uint32             `bson:"given_by" json:"given_by"`
	GivenTo         uint32             `bson:"given_to" json:"given_to"`
}

type CardTransactionRepository interface {
	AddCardTransaction(ctx context.Context, transaction CardTransaction) (uint32, *helpers.CustomError)
	GetCardTransactionsByCardNumber(ctx context.Context, cardNumber string) ([]CardTransaction, *helpers.CustomError)
	GetCardTransactionsByUserId(ctx context.Context, userId uint32) ([]CardTransaction, *helpers.CustomError)
}

type mongoCardTransactionRepo struct {
	collection *mongo.Collection
}

func NewCardTransactionRepository(col *mongo.Collection) CardTransactionRepository {
	return &mongoCardTransactionRepo{collection: col}
}

func (repo *mongoCardTransactionRepo) AddCardTransaction(ctx context.Context, transaction CardTransaction) (uint32, *helpers.CustomError) {
	if transaction.TransactionGuid == 0 {
		nextGuid, err := GetNextSequence(ctx, repo.collection.Database(), "transactionGuids")
		if err != nil {
			return 0, helpers.System("Failed to generate transaction GUID: " + err.Error())
		}
		transaction.TransactionGuid = uint32(nextGuid)
	}
	nextId, err := GetNextSequence(ctx, repo.collection.Database(), "cardTransactions")
	if err != nil {
		return 0, helpers.System("Failed to generate transaction ID: " + err.Error())
	}
	transaction.TransactionId = uint32(nextId)
	_, err = repo.collection.InsertOne(ctx, transaction)
	if err != nil {
		return 0, helpers.System("Failed to add card transaction: " + err.Error())
	}
	return transaction.TransactionId, nil
}

func (repo *mongoCardTransactionRepo) GetCardTransactionsByCardNumber(ctx context.Context, cardNumber string) ([]CardTransaction, *helpers.CustomError) {
	var transactions []CardTransaction
	cursor, err := repo.collection.Find(ctx, bson.M{"card_number": cardNumber})
	if err != nil {
		return nil, helpers.System("Failed to get card transactions by card number: " + err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var transaction CardTransaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, helpers.System("Failed to decode card transaction: " + err.Error())
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, helpers.System("Cursor iteration error: " + err.Error())
	}

	return transactions, nil
}

func (repo *mongoCardTransactionRepo) GetCardTransactionsByUserId(ctx context.Context, userId uint32) ([]CardTransaction, *helpers.CustomError) {
	var transactions []CardTransaction
	cursor, err := repo.collection.Find(ctx, bson.M{"$or": []bson.M{
		{"given_by": userId},
		{"given_to": userId},
	}})
	if err != nil {
		return nil, helpers.System("Failed to get card transactions by user ID: " + err.Error())
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var transaction CardTransaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, helpers.System("Failed to decode card transaction: " + err.Error())
		}
		transactions = append(transactions, transaction)
	}
	if err := cursor.Err(); err != nil {
		return nil, helpers.System("Cursor iteration error: " + err.Error())
	}
	return transactions, nil
}

func (repo *mongoCardTransactionRepo) GetCardTransactionsByTransactionGuid(ctx context.Context, transactionGuid uint32) ([]CardTransaction, *helpers.CustomError) {
	var transactions []CardTransaction
	cursor, err := repo.collection.Find(ctx, bson.M{"transaction_guid": transactionGuid})
	if err != nil {
		return nil, helpers.System("Failed to get card transactions by transaction GUID: " + err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var transaction CardTransaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, helpers.System("Failed to decode card transaction: " + err.Error())
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, helpers.System("Cursor iteration error: " + err.Error())
	}

	return transactions, nil
}
