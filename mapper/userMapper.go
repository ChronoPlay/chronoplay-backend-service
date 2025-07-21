package mapper

import (
	"github.com/gin-gonic/gin"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/models"
)

func DecodeRegisterUserRequest(r *gin.Context) (user model.User, err *helpers.CustomEror) {
	jerr := r.BindJSON(&user)
	if jerr != nil {
		return model.User{}, helpers.BadRequest(jerr.Error())
	}
	return user, nil
}
