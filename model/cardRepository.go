package model

import (
	"context"
	"fmt"
	"log"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Card struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number      string             `bson:"number" json:"number"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Occupied    uint32             `bson:"occupied" json:"occupied"`
	Total       uint32             `bson:"total" json:"total"`
	Available   uint32             `bson:"available" json:"available"`
	Owners      []uint32           `bson:"owners" json:"owners"`
	Creator     uint32             `bson:"creator" json:"creator"`
	ImageUrl    string             `bson:"image_url" json:"image_url"`
}

type GetCardsRequest struct {
	Number  string   `json:"number"`
	Numbers []string `json:"numbers"`
}

type CardRepository interface {
	GetCollection() *mongo.Collection
	AddCard(ctx context.Context, card Card) *helpers.CustomError
	GetAllCards(ctx context.Context) ([]Card, *helpers.CustomError)
	GetCardByNumber(ctx context.Context, cardNumber string) (*Card, *helpers.CustomError)
	UpdateCard(ctx context.Context, card Card) *helpers.CustomError
	GetOwnersByCardNumber(ctx context.Context, cardNumber string) ([]uint32, *helpers.CustomError)
	GetCards(ctx context.Context, req GetCardsRequest) ([]Card, *helpers.CustomError)
	UpdateCards(ctx context.Context, cards []Card) *helpers.CustomError
}

type mongoCardRepo struct {
	collection *mongo.Collection
}

func NewCardRepository(col *mongo.Collection) CardRepository {
	return &mongoCardRepo{collection: col}
}

func (r *mongoCardRepo) GetCollection() *mongo.Collection {
	return r.collection
}

func (repo *mongoCardRepo) AddCard(ctx context.Context, card Card) *helpers.CustomError {
	_, err := repo.collection.InsertOne(ctx, card)
	if err != nil {
		return helpers.System(fmt.Sprintf("%s: %s", err.Error(), "Failed to add card"))
	}
	return nil
}

func (repo *mongoCardRepo) GetAllCards(ctx context.Context) ([]Card, *helpers.CustomError) {
	var cards []Card
	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, helpers.System(fmt.Sprintf("%s: %s", err.Error(), "Failed to get all cards"))
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var card Card
		if err := cursor.Decode(&card); err != nil {
			return nil, helpers.System(fmt.Sprintf("%s: %s", err.Error(), "Failed to decode card"))
		}
		cards = append(cards, card)
	}

	if err := cursor.Err(); err != nil {
		return nil, helpers.System(fmt.Sprintf("%s: %s", err.Error(), "Failed to iterate cursor"))
	}

	return cards, nil
}

func (repo *mongoCardRepo) GetCardByNumber(ctx context.Context, cardNumber string) (*Card, *helpers.CustomError) {
	var card Card
	err := repo.collection.FindOne(ctx, bson.M{"number": cardNumber}).Decode(&card)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, helpers.NotFound("card not found")
		}
		return nil, helpers.System(fmt.Sprintf("%s: %s", err.Error(), "Failed to find card by number"))
	}
	return &card, nil
}

func (repo *mongoCardRepo) UpdateCard(ctx context.Context, card Card) *helpers.CustomError {
	card.ID = primitive.NewObjectID() // Ensure ID is set for update
	updateData, err := bson.Marshal(card)
	if err != nil {
		return helpers.System(fmt.Sprintf("%s: %s", err.Error(), "Failed to marshal card for update"))
	}

	_, err = repo.collection.UpdateOne(ctx, bson.M{"_id": card.ID}, bson.M{"$set": updateData})
	if err != nil {
		return helpers.System(fmt.Sprintf("%s: %s", err.Error(), "Failed to update card"))
	}
	return nil
}

func (repo *mongoCardRepo) UpdateCards(ctx context.Context, cards []Card) *helpers.CustomError {
	if len(cards) == 0 {
		return helpers.BadRequest("No cards to update")
	}

	for _, card := range cards {
		card.ID = primitive.NewObjectID() // Ensure ID is set for update
		updateData, err := bson.Marshal(card)
		if err != nil {
			return helpers.System(fmt.Sprintf("%s: %s", err.Error(), "Failed to marshal card for update"))
		}

		_, err = repo.collection.UpdateOne(ctx, bson.M{"_id": card.ID}, bson.M{"$set": updateData})
		if err != nil {
			return helpers.System(fmt.Sprintf("%s: %s", err.Error(), "Failed to update card"))
		}
	}
	return nil
}

func (repo *mongoCardRepo) GetOwnersByCardNumber(ctx context.Context, cardNumber string) ([]uint32, *helpers.CustomError) {
	var card Card
	err := repo.collection.FindOne(ctx, bson.M{"number": cardNumber}).Decode(&card)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, helpers.NotFound("card not found")
		}
		return nil, helpers.System(fmt.Sprintf("%s: %s", err.Error(), "Failed to find card by number"))
	}
	if len(card.Owners) == 0 {
		return nil, helpers.NotFound("no owners found for this card")
	}
	return card.Owners, nil
}

func (repo *mongoCardRepo) GetCards(ctx context.Context, req GetCardsRequest) ([]Card, *helpers.CustomError) {
	cards := []Card{}
	isValid := false
	conditions := []bson.M{}
	if len(req.Numbers) > 0 {
		conditions = append(conditions, bson.M{
			"number": bson.M{"$in": req.Numbers},
		})
		isValid = true
	} else if req.Number != "" {
		conditions = append(conditions, bson.M{
			"number": req.Number,
		})
		isValid = true
	}
	if isValid {
		filter := bson.M{}
		if len(conditions) > 0 {
			filter["$and"] = conditions
		}
		cursor, err := repo.collection.Find(ctx, filter)
		if err != nil {
			return cards, helpers.System(err.Error())
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var card Card
			if err := cursor.Decode(&card); err != nil {
				log.Println("Decode error:", err)
				continue
			}
			cards = append(cards, card)
		}

		if err := cursor.Err(); err != nil {
			return cards, helpers.System(err.Error())
		}
	}
	return cards, nil
}
