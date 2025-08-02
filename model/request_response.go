package model

type VerifyUserRequest struct {
	Email string `bson:"email" json:"email"`
}

type LoginUserRequest struct {
	Email       string `bson:"email" json:"email"`
	PhoneNumber string `bson:"phone_number" json:"phone_number"`
	Password    string `bson:"password" json:"password"`
}

type LoginUserResponse struct {
	JwtToken string `json:"jwt_token"`
}
