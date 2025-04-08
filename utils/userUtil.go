package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/models"
)

func ValidateUser(user model.User) (err *helpers.CustomEror) {
	if len(strings.TrimSpace(user.UserName)) == 0 {
		return err.BadRequest("username is required")
	}
	if len(strings.TrimSpace(user.Name)) == 0 {
		return err.BadRequest("name is required")
	}
	if len(strings.TrimSpace(user.Email)) == 0 {
		return err.BadRequest("email is required")
	}
	if !isValidEmail(user.Email) {
		return err.BadRequest("invalid email format")
	}
	if len(strings.TrimSpace(user.Password)) < 6 {
		return err.BadRequest("password must be at least 6 characters")
	}
	if len(strings.TrimSpace(user.PhoneNumber)) < 10 {
		return err.BadRequest("phone number is too short")
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func SendEmailToUser(req dto.EmailVerificationRequest) (err *helpers.CustomEror) {
	if !isValidEmail(req.Email) {
		return err.BadRequest("invalid email format")
	}

	subject := "Email verification for chronoplay"
	body := GenerateVerificationEmailBody(req.Link, req.UserName)

	err = SendEmail([]string{req.Email}, subject, body)
	if err != nil {
		return err
	}

	return nil
}

func GenerateVerificationEmailBody(verificationLink string, userName string) string {
	return fmt.Sprintf(`
Hello %s,

Thanks for registering with us!

Please verify your email address by clicking the button below:

🔗 Verify Email: %s

If you didn't request this, you can safely ignore this email.

Best regards,  
The ChronoPlay Team
`, userName, verificationLink)
}

func GenrateEmailVerificationLink(email string) string {
	baseUrl := os.Getenv("BASE_URL")
	link := fmt.Sprintf(`%s/verifyEmail?email=%s`, baseUrl, email)
	return link
}

func HashPassword(password string) (hashedPassword string, err *helpers.CustomEror) {
	bytes, gerr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if gerr != nil {
		return "", err.System("error while hashing password: " + gerr.Error())
	}
	return string(bytes), nil
}
