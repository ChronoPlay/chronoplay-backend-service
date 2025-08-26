package model

import (
	"context"
	"log"
	"time"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CashTransaction struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TransactionId   uint32             `bson:"transaction_id" json:"transaction_id"`
	TransactionGuid uint32             `bson:"transaction_guid" json:"transaction_guid"`
	Amount          float32            `bson:"amount" json:"amount"`
	GivenBy         uint32             `bson:"given_by" json:"given_by"`
	GivenTo         uint32             `bson:"given_to" json:"given_to"`
	Status          string             `bson:"status" json:"status"`
	CreatedAt       primitive.DateTime `bson:"created_at" json:"created_at"`
	CreatedBy       uint32             `bson:"created_by" json:"created_by"`
	UpdatedBy       uint32             `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
	UpdatedAt       primitive.DateTime `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type CashTransactionRepository interface {
	GetCollection() *mongo.Collection
	AddCashTransaction(ctx context.Context, transaction CashTransaction) (uint32, *helpers.CustomError)
	GetCashTransactionsByUserId(ctx context.Context, userId uint32) ([]CashTransaction, *helpers.CustomError)
	GetCashTransactionsToUserId(ctx context.Context, userId uint32) ([]CashTransaction, *helpers.CustomError)
	GetCashTransactionsByTransactionId(ctx context.Context, transactionId uint32) ([]CashTransaction, *helpers.CustomError)
	GetCashTransactionsByTransactionGuid(ctx context.Context, transactionGuid uint32) ([]CashTransaction, *helpers.CustomError)
	UpdateCashTransactions(ctx context.Context, transactions []CashTransaction) *helpers.CustomError
}

type mongoCashTransactionRepo struct {
	collection *mongo.Collection
}

func NewCashTransactionRepository(col *mongo.Collection) CashTransactionRepository {
	return &mongoCashTransactionRepo{collection: col}
}

func (repo *mongoCashTransactionRepo) GetCollection() *mongo.Collection {
	return repo.collection
}

func (repo *mongoCashTransactionRepo) AddCashTransaction(ctx context.Context, transaction CashTransaction) (uint32, *helpers.CustomError) {
	if transaction.TransactionGuid == 0 {
		nextGuid, err := GetNextSequence(ctx, repo.collection.Database(), "transactionGuids")
		if err != nil {
			return 0, helpers.System("Failed to generate transaction GUID: " + err.Error())
		}
		transaction.TransactionGuid = uint32(nextGuid)
	}
	log.Printf("Generated Transaction GUID: %d\n", transaction.TransactionGuid)
	nextId, err := GetNextSequence(ctx, repo.collection.Database(), "cashTransactions")
	if err != nil {
		return 0, helpers.System("Failed to generate transaction ID: " + err.Error())
	}
	transaction.TransactionId = uint32(nextId)
	transaction.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err = repo.collection.InsertOne(ctx, transaction)
	if err != nil {
		return 0, helpers.System("Failed to add cash transaction: " + err.Error())
	}
	return transaction.TransactionId, nil
}

func (repo *mongoCashTransactionRepo) GetCashTransactionsByUserId(ctx context.Context, userId uint32) ([]CashTransaction, *helpers.CustomError) {
	var transactions []CashTransaction
	cursor, err := repo.collection.Find(ctx, bson.M{"given_by": userId})
	if err != nil {
		return nil, helpers.System("Failed to get cash transactions by user ID: " + err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var transaction CashTransaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, helpers.System("Failed to decode cash transaction: " + err.Error())
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, helpers.System("Failed to iterate cursor: " + err.Error())
	}

	return transactions, nil
}

func (repo *mongoCashTransactionRepo) GetCashTransactionsToUserId(ctx context.Context, userId uint32) ([]CashTransaction, *helpers.CustomError) {
	var transactions []CashTransaction
	cursor, err := repo.collection.Find(ctx, bson.M{"given_to": userId})
	if err != nil {
		return nil, helpers.System("Failed to get cash transactions to user ID: " + err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var transaction CashTransaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, helpers.System("Failed to decode cash transaction: " + err.Error())
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, helpers.System("Failed to iterate cursor: " + err.Error())
	}

	return transactions, nil
}

func (repo *mongoCashTransactionRepo) GetCashTransactionsByTransactionId(ctx context.Context, transactionId uint32) ([]CashTransaction, *helpers.CustomError) {
	var transactions []CashTransaction
	cursor, err := repo.collection.Find(ctx, bson.M{"transaction_id": transactionId})
	if err != nil {
		return nil, helpers.System("Failed to get cash transactions by transaction ID: " + err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var transaction CashTransaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, helpers.System("Failed to decode cash transaction: " + err.Error())
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, helpers.System("Failed to iterate cursor: " + err.Error())
	}

	return transactions, nil
}

func (repo *mongoCashTransactionRepo) GetCashTransactionsByTransactionGuid(ctx context.Context, transactionGuid uint32) ([]CashTransaction, *helpers.CustomError) {
	var transactions []CashTransaction
	cursor, err := repo.collection.Find(ctx, bson.M{"transaction_guid": transactionGuid})
	if err != nil {
		return nil, helpers.System("Failed to get cash transactions by transaction GUID: " + err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var transaction CashTransaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, helpers.System("Failed to decode cash transaction: " + err.Error())
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, helpers.System("Failed to iterate cursor: " + err.Error())
	}

	return transactions, nil
}

func (repo *mongoCashTransactionRepo) UpdateCashTransactions(ctx context.Context, transactions []CashTransaction) *helpers.CustomError {
	for _, transaction := range transactions {
		filter := bson.M{"transaction_id": transaction.TransactionId}
		update := bson.M{
			"$set": bson.M{
				"amount":           transaction.Amount,
				"given_by":         transaction.GivenBy,
				"given_to":         transaction.GivenTo,
				"status":           transaction.Status,
				"updated_at":       primitive.NewDateTimeFromTime(time.Now()),
				"updated_by":       transaction.UpdatedBy,
				"transaction_guid": transaction.TransactionGuid,
			},
		}
		_, err := repo.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return helpers.System("Failed to update cash transaction: " + err.Error())
		}
	}
	return nil
}
