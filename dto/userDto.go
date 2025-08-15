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
