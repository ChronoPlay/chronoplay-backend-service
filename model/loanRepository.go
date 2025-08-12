package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Loan struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Amount              float32            `bson:"amount" json:"amount"`
	LoanedBy            uint32             `bson:"loaned_by" json:"loaned_by"`
	LoanedTo            uint32             `bson:"loaned_to" json:"loaned_to"`
	Rate                float32            `bson:"rate" json:"rate"` // compound interest applied after each day
	InterestAccumulated float32            `bson:"interest_accumulated" json:"interest_accumulated"`
	PaidAmount          float32            `bson:"paid_amount" json:"paid_amount"`
	IsPaid              bool               `bson:"is_paid" json:"is_paid"`
}
