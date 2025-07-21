package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	UserName     string             `bson:"username"`
	PhoneNumber  string             `bson:"phone_number"`
	Cash         uint32             `bson:"cash"`
	Bronze       uint32             `bson:"bronze"`
	Silver       uint32             `bson:"silver"`
	Gold         uint32             `bson:"gold"`
	IsAuthorized bool               `bson:"is_authorized"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

type UserRepository interface {
	FindByUserName(ctx context.Context, username string) (*User, *helpers.CustomEror)
	RegisterUser(sessCtx mongo.SessionContext, user User) *helpers.CustomEror
	GetCollection() *mongo.Collection
}

type mongoUserRepo struct {
	collection *mongo.Collection
}

func NewUserRepository(col *mongo.Collection) UserRepository {
	return &mongoUserRepo{collection: col}
}

func (repo *mongoUserRepo) FindByUserName(ctx context.Context, userName string) (*User, *helpers.CustomEror) {
	var user User
	err := repo.collection.FindOne(ctx, bson.M{"username": userName}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, helpers.NotFound("user not found")
		}
		return nil, helpers.System(err.Error())
	}
	return &user, nil
}

func (repo *mongoUserRepo) RegisterUser(sessCtx mongo.SessionContext, user User) *helpers.CustomEror {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := repo.collection.InsertOne(sessCtx, user)
	if err != nil {
		return helpers.System(err.Error())
	}
	return nil
}

func (r *mongoUserRepo) GetCollection() *mongo.Collection {
	return r.collection
}
