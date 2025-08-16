package dto

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
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
	UserName string `bson:"user_name" json:"user_name"`
}

type AddFriendRequest struct {
	UserID   uint32 `bson:"user_id" json:"user_id"`
	FriendID uint32 `bson:"friend_id" json:"friend_id"`
}

type AddFriendResponse struct {
	Message string `bson:"message" json:"message"`
}
