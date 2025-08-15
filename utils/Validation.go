package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/model"
)

func ValidateUser(user model.User) (err *helpers.CustomError) {
	if len(strings.TrimSpace(user.UserName)) == 0 {
		return helpers.BadRequest("username is required")
	}
	if len(strings.TrimSpace(user.Name)) == 0 {
		return helpers.BadRequest("name is required")
	}
	if len(strings.TrimSpace(user.Email)) == 0 {
		return helpers.BadRequest("email is required")
	}
	if !IsValidEmail(user.Email) {
		return helpers.BadRequest("invalid email format")
	}
	if len(strings.TrimSpace(user.Password)) < 6 {
		return helpers.BadRequest("password must be at least 6 characters")
	}
	if len(strings.TrimSpace(user.PhoneNumber)) < 10 {
		return helpers.BadRequest("phone number is too short")
	}
	return nil
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)
	return re.MatchString(email)
}

func SendEmailToUser(req dto.EmailVerificationRequest) (err *helpers.CustomError) {
	if !IsValidEmail(req.Email) {
		return helpers.BadRequest("invalid email format")
	}

	subject := "Email verification for chronoplay"
	body := GenerateVerificationEmailBody(req.Link, req.UserName)

	fmt.Println("Sending mail now....")
	err = SendEmail([]string{req.Email}, subject, body)
	if err != nil {
		return err
	}
	fmt.Println("Mail sent for verification suceesfully...")

	return nil
}

func GenerateVerificationEmailBody(verificationLink string, userName string) string {
	return fmt.Sprintf(`
<html>
<body>
<p>Hello %s,</p>

<p>Thanks for registering with us!</p>

<p>Please verify your email address by clicking the link below:</p>

<p><a href="%s" style="padding: 10px 20px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 4px;">Verify Email</a></p>

<p>If you didn't request this, you can safely ignore this email.</p>

<p>Best regards,<br>
The ChronoPlay Team</p>
</body>
</html>
`, userName, verificationLink)
}

func GenrateEmailVerificationLink(email string) string {
	baseUrl := os.Getenv("BASE_URL")
	link := fmt.Sprintf(`%s/auth/verify?email=%s`, baseUrl, email)
	return link
}

func HashPassword(password string) (hashedPassword string, err *helpers.CustomError) {
	bytes, gerr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if gerr != nil {
		return "", helpers.System("error while hashing password: " + gerr.Error())
	}
	return string(bytes), nil
}

func CheckPasswordHash(password string, hashedPassword string) (err *helpers.CustomError) {
	println("Checking password hash...")
	println("password: ", password,
		"\nhashedPassword: ", hashedPassword)
	berr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if berr != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil
		}
		return helpers.System("error while comparing password hash: " + err.Error())
	}
	return nil
}

func GenerateJwtToken(userId uint32, userType string) (jwtToken string, err *helpers.CustomError) {
	claims := jwt.MapClaims{
		"user_id":   userId,
		"user_type": userType,
		"exp":       time.Now().Add(time.Hour * 1).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	signedToken, jerr := token.SignedString(jwtSecret)
	if jerr != nil {
		return "", helpers.System("error while signing JWT token: " + jerr.Error())
	}

	return signedToken, nil
}

func ParseJwtToken(tokenString string) (userId uint32, userType string, err *helpers.CustomError) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	token, terr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, helpers.System("unexpected signing method: " + token.Header["alg"].(string))
		}
		return jwtSecret, nil
	})

	if terr != nil || !token.Valid {
		return 0, userType, helpers.Unauthorized("invalid JWT token: " + terr.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil || claims["user_type"] == nil {
		return 0, userType, helpers.Unauthorized("invalid JWT claims")
	}

	userId = uint32(claims["user_id"].(float64))
	//userType = claims["user_type"].(string)
	return userId, userType, nil
}

func IsAdmin(userType string) bool {
	return strings.EqualFold(userType, model.USER_TYPE_ADMIN)
}

func ValidateAddCardRequest(req dto.AddCardRequest) (err *helpers.CustomError) {
	if len(strings.TrimSpace(req.CardNumber)) == 0 {
		return helpers.BadRequest("card number is required")
	}
	if len(strings.TrimSpace(req.CardDescription)) == 0 {
		return helpers.BadRequest("card description is required")
	}
	if req.TotalCards == 0 {
		return helpers.BadRequest("total cards must be greater than zero")
	}
	if req.UserId == 0 {
		return helpers.BadRequest("user ID is required")
	}
	if len(strings.TrimSpace(req.UserType)) == 0 {
		return helpers.BadRequest("user type is required")
	}
	if req.Image == nil {
		return helpers.BadRequest("image is required")
	}
	if !IsAdmin(req.UserType) {
		return helpers.Unauthorized("only admin can add cards")
	}
	return nil
}

func ValidateGetCardRequest(req dto.GetCardRequest) (err *helpers.CustomError) {
	if len(strings.TrimSpace(req.CardNumber)) == 0 {
		return helpers.BadRequest("card number is required")
	}
	return nil
}

func ValidateTransferCashRequest(req dto.TransferCashRequest) (err *helpers.CustomError) {
	if req.Amount <= 0 {
		return helpers.BadRequest("amount must be greater than zero")
	}
	if req.GivenBy == 0 {
		return helpers.BadRequest("given by user ID is required")
	}
	if req.GivenTo == 0 {
		return helpers.BadRequest("given to user ID is required")
	}
	if req.GivenBy == req.GivenTo {
		return helpers.BadRequest("given by and given to user IDs cannot be the same")
	}
	if req.UserId == 0 {
		return helpers.BadRequest("user ID is required")
	}
	if req.Status == model.TRANSACTION_STATUS_SUCCESS && req.UserId != req.GivenBy {
		return helpers.BadRequest("only the user who is giving the amount can mark it as successful")
	}
	return nil
}

func IsValidTransactionStatus(status string) bool {
	for _, validStatus := range model.ValidTransactionStatuses {
		if strings.EqualFold(status, validStatus) {
			return true
		}
	}
	return false
}

func ValidateTransferCardsRequest(req dto.TransferCardRequest) (err *helpers.CustomError) {
	if len(req.Cards) == 0 {
		return helpers.BadRequest("at least one card is required for transfer")
	}
	if req.GivenBy == 0 {
		return helpers.BadRequest("given by user ID is required")
	}
	if req.GivenTo == 0 {
		return helpers.BadRequest("given to user ID is required")
	}
	if req.GivenBy == req.GivenTo {
		return helpers.BadRequest("given by and given to user IDs cannot be the same")
	}
	if req.Status == model.TRANSACTION_STATUS_SUCCESS && req.UserId != req.GivenBy {
		return helpers.BadRequest("only the user who is giving the cards can mark it as successful")
	}
	return nil
}
