package dto

type EmailVerificationRequest struct {
	Email    string
	UserName string
	Link     string
}
