package mapper

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

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
	user.UserType = model.USER_TYPE_USER

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

func MapCardsToResponse(cards []model.Card) []dto.CardResponse {
	var cardResponses []dto.CardResponse
	for _, card := range cards {
		cardResponses = append(cardResponses, dto.CardResponse{
			Number:      card.Number,
			Occupied:    card.Occupied,
			Image:       card.ImageUrl,
			Description: card.Description,
			Rarity:      card.Rarity,
			Name:        card.Name,
		})
	}
	return cardResponses
}

func EncodeGetUserByIdResponse(req dto.GetUserResponse) (res dto.GetUserByIdResponse) {
	res.Email = req.Email
	res.Name = req.Name
	res.UserName = req.UserName
	res.Cash = req.Cash
	res.UserType = req.UserType
	res.Cards = req.Cards
	return res
}

func DecodeAddFriendRequest(curUserId interface{}, friendUserId string) (*dto.AddFriendRequest, error) {
	uid, ok := curUserId.(uint32)
	if !ok {
		return nil, errors.New("invalid current user id")
	}

	fid, err := strconv.ParseInt(friendUserId, 10, 32)
	if err != nil {
		return nil, err
	}

	return &dto.AddFriendRequest{
		UserID:   uid,
		FriendID: uint32(fid),
	}, nil
}
