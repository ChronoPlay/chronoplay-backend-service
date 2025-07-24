package mapper

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/models"
)

func DecodeRegisterUserRequest(r *gin.Context) (user model.User, err *helpers.CustomEror) {
	if err := r.ShouldBindJSON(&user); err != nil {
		return model.User{}, helpers.BadRequest("Invalid request: " + err.Error())
	}

	// Log the final parsed user
	if bytes, _ := json.Marshal(user); bytes != nil {
		log.Println("Parsed user:", string(bytes))
	}

	return user, nil
}

func DecodeVerifyUserRequest(r *gin.Context) (req model.VerifyUserRequest, err *helpers.CustomEror) {
	if err := r.ShouldBindJSON(&req); err != nil {
		return model.VerifyUserRequest{}, helpers.BadRequest("Invalid request: " + err.Error())
	}

	// Log the final parsed user
	if bytes, _ := json.Marshal(req); bytes != nil {
		log.Println("Parsed user:", string(bytes))
	}

	return req, nil
}
