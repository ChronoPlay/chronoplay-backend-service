package dto

import "github.com/ChronoPlay/chronoplay-backend-service/model"

type EmailVerificationRequest struct {
	Email    string
	UserName string
	Link     string
}

type VerifyUserRequest struct {
	Email string `bson:"email" json:"email"`
}

type LoginUserRequest struct {
	Email      string `bson:"email" json:"email"`
	UserName   string `bson:"user_name" json:"user_name"`
	Password   string `bson:"password" json:"password"`
	Identifier string `bson:"identifier" json:"identifier"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

type GetUserResponse struct {
	Name        string               `json:"name"`
	Email       string               `json:"email"`
	UserName    string               `json:"user_name"`
	Cash        float32              `json:"cash"`
	FriendIds   []uint32             `json:"friend_ids"`
	PhoneNumber string               `bson:"phone_number" json:"phone_number"`
	Cards       []model.CardOccupied `bson:"cards" json:"cards"`
	UserType    string               `json:"user_type"`
}

type GetUserByIdResponse struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	UserName string  `json:"user_name"`
	Cash     float32 `json:"cash"`
	UserType string  `json:"user_type"`
}

type AddFriendRequest struct {
	UserID   uint32 `bson:"user_id" json:"user_id"`
	FriendID uint32 `bson:"friend_id" json:"friend_id"`
}

type AddFriendResponse struct {
	Message string `bson:"message" json:"message"`
}

type GetFriendsRequest struct {
	UserID uint32 `bson:"user_id" json:"user_id"`
}

type Friend struct {
	UserID   uint32 `bson:"user_id" json:"user_id"`
	UserName string `bson:"user_name" json:"user_name"`
	Email    string `bson:"email" json:"email"`
}
type GetFriendsResponse struct {
	Friends []Friend `bson:"friends" json:"friends"`
}
