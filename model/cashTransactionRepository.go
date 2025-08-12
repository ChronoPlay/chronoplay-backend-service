package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CashTransaction struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TransactionId uint32             `bson:"transaction_id" json:"transaction_id"`
	Amount        float32            `bson:"amount" json:"amount"`
	GivenBy       uint32             `bson:"given_by" json:"given_by"`
	GivenTo       uint32             `bson:"given_to" json:"given_to"`
}
