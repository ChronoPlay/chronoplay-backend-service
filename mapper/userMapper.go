package mapper

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/model"
	utils "github.com/ChronoPlay/chronoplay-backend-service/utils"
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

func DecodeVerifyUserRequest(r *gin.Context) (req dto.VerifyUserRequest, err *helpers.CustomError) {
	// Get the "email" from query parameters
	email := r.Query("email")
	if email == "" {
		return dto.VerifyUserRequest{}, helpers.BadRequest("Missing email in query parameters")
	}

	req.Email = email

	log.Println("Parsed user with email:", req.Email)

	return req, nil
}

func DecodeLoginUserRequest(r *gin.Context) (req dto.LoginUserRequest, err *helpers.CustomError) {

	if err := r.ShouldBindJSON(&req); err != nil {
		return dto.LoginUserRequest{}, helpers.BadRequest("Invalid request: " + err.Error())
	}
	if req.Identifier != "" {
		if utils.IsValidEmail(req.Identifier) {
			req.Email = req.Identifier
		} else {
			req.UserName = req.Identifier
		}
	}

	if req.Email == "" && req.UserName == "" {
		return dto.LoginUserRequest{}, helpers.BadRequest("Missing email or user_name in request body")
	}

	if req.Password == "" {
		return dto.LoginUserRequest{}, helpers.BadRequest("Missing password in request body")
	}

	log.Println("Parsed login request with email:", req.Email, "and user_name:", req.UserName)
	return req, nil
}
