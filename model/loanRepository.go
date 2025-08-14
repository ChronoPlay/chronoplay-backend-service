package model

import (
	"context"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Loan struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LoanId              uint32             `bson:"loan_id" json:"loan_id"`
	Amount              float32            `bson:"amount" json:"amount"`
	LoanedBy            uint32             `bson:"loaned_by" json:"loaned_by"`
	LoanedTo            uint32             `bson:"loaned_to" json:"loaned_to"`
	Rate                float32            `bson:"rate" json:"rate"` // compound interest applied after each day
	InterestAccumulated float32            `bson:"interest_accumulated" json:"interest_accumulated"`
	PaidAmount          float32            `bson:"paid_amount" json:"paid_amount"`
	IsPaid              bool               `bson:"is_paid" json:"is_paid"`
}

type LoanRepository interface {
	AddLoan(ctx context.Context, loan Loan) (uint32, *helpers.CustomError)
	GetLoansByUserId(ctx context.Context, userId uint32) ([]Loan, *helpers.CustomError)
	GetLoanByLoanId(ctx context.Context, loanId uint32) (*Loan, *helpers.CustomError)
	UpdateLoan(ctx context.Context, loan Loan) *helpers.CustomError
}

type mongoLoanRepo struct {
	collection *mongo.Collection
}

func NewLoanRepository(col *mongo.Collection) LoanRepository {
	return &mongoLoanRepo{collection: col}
}

func (repo *mongoLoanRepo) AddLoan(ctx context.Context, loan Loan) (uint32, *helpers.CustomError) {
	nextId, err := GetNextSequence(ctx, repo.collection.Database(), "loans")
	if err != nil {
		return 0, helpers.System("Failed to generate loan ID: " + err.Error())
	}
	loan.LoanId = uint32(nextId)
	_, err = repo.collection.InsertOne(ctx, loan)
	if err != nil {
		return 0, helpers.System("Failed to add loan: " + err.Error())
	}
	return loan.LoanId, nil
}

func (repo *mongoLoanRepo) GetLoansByUserId(ctx context.Context, userId uint32) ([]Loan, *helpers.CustomError) {
	var loans []Loan
	cursor, err := repo.collection.Find(ctx, bson.M{"$or": []bson.M{
		{"loaned_by": userId},
		{"loaned_to": userId},
	}})
	if err != nil {
		return nil, helpers.System("Failed to get loans by user ID: " + err.Error())
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var loan Loan
		if err := cursor.Decode(&loan); err != nil {
			return nil, helpers.System("Failed to decode loan: " + err.Error())
		}
		loans = append(loans, loan)
	}
	if err := cursor.Err(); err != nil {
		return nil, helpers.System("Failed to iterate cursor: " + err.Error())
	}
	return loans, nil
}

func (repo *mongoLoanRepo) GetLoanByLoanId(ctx context.Context, loanId uint32) (*Loan, *helpers.CustomError) {
	var loan Loan
	err := repo.collection.FindOne(ctx, bson.M{"loan_id": loanId}).Decode(&loan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, helpers.NotFound("loan not found")
		}
		return nil, helpers.System("Failed to find loan by ID: " + err.Error())
	}
	return &loan, nil
}

func (repo *mongoLoanRepo) UpdateLoan(ctx context.Context, loan Loan) *helpers.CustomError {
	result, err := repo.collection.UpdateOne(ctx, bson.M{"loan_id": loan.LoanId}, bson.M{"$set": loan})
	if err != nil {
		return helpers.System("Failed to update loan: " + err.Error())
	}
	if result.MatchedCount == 0 {
		return helpers.NotFound("loan not found")
	}
	return nil
}
