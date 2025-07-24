package model

type VerifyUserRequest struct {
	UserId uint32 `bson:"user_id" json:"user_id"`
}
