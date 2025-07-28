package model

type VerifyUserRequest struct {
	Email string `bson:"email" json:"email"`
}
