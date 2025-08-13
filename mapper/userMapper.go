package mapper

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/model"
)

func DecodeRegisterUserRequest(r *gin.Context) (user model.User, err *helpers.CustomError) {
	if err := r.ShouldBindJSON(&user); err != nil {
		return model.User{}, helpers.BadRequest("Invalid request: " + err.Error())
	}

	// Log the final parsed user
	if bytes, _ := json.Marshal(user); bytes != nil {
		log.Println("Parsed user:", string(bytes))
	}

	return user, nil
}

func DecodeVerifyUserRequest(r *gin.Context) (req model.VerifyUserRequest, err *helpers.CustomError) {
	// Get the "email" from query parameters
	email := r.Query("email")
	if email == "" {
		return model.VerifyUserRequest{}, helpers.BadRequest("Missing email in query parameters")
	}

	req.Email = email

	log.Println("Parsed user with email:", req.Email)

	return req, nil
}

func DecodeLoginUser(r *gin.Context) (req model.LoginUserRequest, err *helpers.CustomError) {
	if err := r.ShouldBindJSON(&req); err != nil {
		return model.LoginUserRequest{}, helpers.BadRequest("Invalid request: " + err.Error())
	}

	if req.Email == "" && req.PhoneNumber == "" {
		return model.LoginUserRequest{}, helpers.BadRequest("Missing email or phone_number in request body")
	}

	if req.Password == "" {
		return model.LoginUserRequest{}, helpers.BadRequest("Missing password in request body")
	}

	log.Println("Parsed login request with email:", req.Email, "and phone number:", req.PhoneNumber)
	return req, nil
}
