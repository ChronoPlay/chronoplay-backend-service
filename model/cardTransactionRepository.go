package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CardTransaction struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TransactionId uint32             `bson:"transaction_id" json:"transaction_id"`
	Amount        float32            `bson:"amount" json:"amount"`
	CardNumber    string             `bson:"card_number" json:"card_number"`
	GivenBy       uint32             `bson:"given_by" json:"given_by"`
	GivenTo       uint32             `bson:"given_to" json:"given_to"`
}
