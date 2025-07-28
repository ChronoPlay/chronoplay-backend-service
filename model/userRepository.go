package model

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId       uint32             `bson:"user_id" json:"user_id"`
	Name         string             `bson:"name" json:"name"`
	Email        string             `bson:"email" json:"email"`
	Password     string             `bson:"password" json:"password"`
	UserName     string             `bson:"username" json:"username"`
	PhoneNumber  string             `bson:"phone_number" json:"phone_number"`
	Cash         uint32             `bson:"cash" json:"cash"`
	Bronze       uint32             `bson:"bronze" json:"bronze"`
	Silver       uint32             `bson:"silver" json:"silver"`
	Gold         uint32             `bson:"gold" json:"gold"`
	IsAuthorized bool               `bson:"is_authorized" json:"is_authorized"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type UserRepository interface {
	FindByUserName(ctx context.Context, username string) (*User, *helpers.CustomEror)
	RegisterUser(sessCtx mongo.SessionContext, user User) (uint32, *helpers.CustomEror)
	GetCollection() *mongo.Collection
	GetUsers(ctx context.Context, req User) ([]User, *helpers.CustomEror)
	UpdateUser(ctx context.Context, user User) *helpers.CustomEror
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

func (repo *mongoUserRepo) RegisterUser(sessCtx mongo.SessionContext, user User) (uint32, *helpers.CustomEror) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	userId, err := GetNextSequence(sessCtx, repo.collection.Database(), COUNTER_ID_USER_ID)
	if err != nil {
		return 0, helpers.System(err.Error())
	}
	user.UserId = uint32(userId)
	_, err = repo.collection.InsertOne(sessCtx, user)
	if err != nil {
		return 0, helpers.System(err.Error())
	}
	return user.UserId, nil
}

func (r *mongoUserRepo) GetCollection() *mongo.Collection {
	return r.collection
}

func (r *mongoUserRepo) UpdateUser(ctx context.Context, user User) *helpers.CustomEror {
	user.UpdatedAt = time.Now()

	// Convert struct to bson.M
	updateData, err := bson.Marshal(user)
	if err != nil {
		return helpers.System("failed to marshal user: " + err.Error())
	}

	var updateDoc bson.M
	if err := bson.Unmarshal(updateData, &updateDoc); err != nil {
		return helpers.System("failed to unmarshal user: " + err.Error())
	}

	// Wrap in $set
	update := bson.M{
		"$set": updateDoc,
	}

	_, err = r.collection.UpdateByID(ctx, user.ID, update)
	if err != nil {
		return helpers.System(err.Error())
	}
	return nil
}

func (r *mongoUserRepo) GetUsers(ctx context.Context, req User) ([]User, *helpers.CustomEror) {
	users := []User{}
	isValid := false
	conditions := []bson.M{}
	if req.UserId != 0 {
		conditions = append(conditions, bson.M{
			"user_id": req.UserId,
		})
		isValid = true
	}
	if req.ID != primitive.NilObjectID {
		conditions = append(conditions, bson.M{
			"_id": req.ID,
		})
		isValid = true
	}
	if req.UserName != "" {
		conditions = append(conditions, bson.M{
			"user_name": req.UserName,
		})
		isValid = true
	}
	if req.Email != "" {
		conditions = append(conditions, bson.M{
			"email": req.Email,
		})
		isValid = true
	}

	if isValid {
		filter := bson.M{}
		if len(conditions) > 0 {
			filter["$and"] = conditions
		}
		cursor, err := r.collection.Find(ctx, filter)
		if err != nil {
			return users, helpers.System(err.Error())
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var user User
			if err := cursor.Decode(&user); err != nil {
				log.Println("Decode error:", err)
				continue
			}
			users = append(users, user)
		}

		if err := cursor.Err(); err != nil {
			return users, helpers.System(err.Error())
		}
	}
	return users, nil
}
