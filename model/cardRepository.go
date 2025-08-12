package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Card struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number      string             `bson:"card_number" json:"card_number"`
	Description string             `bson:"description" json:"description"`
	Occupied    uint32             `bson:"occupied" json:"occupied"`
	Total       uint32             `bson:"total" json:"total"`
	Available   uint32             `bson:"available" json:"available"`
	Owners      []uint32           `bson:"owners" json:"owners"`
	Creator     uint32             `bson:"creator" json:"creator"`
}
