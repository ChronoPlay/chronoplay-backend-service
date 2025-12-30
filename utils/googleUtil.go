package utils

import (
	"context"
	"errors"

	model "github.com/ChronoPlay/chronoplay-backend-service/model"
	"google.golang.org/api/idtoken"
)

func VerifyGoogleIDToken(ctx context.Context, token, clientID string) (*model.GooglePayload, error) {
	payload, err := idtoken.Validate(ctx, token, clientID)
	if err != nil {
		return nil, err
	}

	if payload.Issuer != "https://accounts.google.com" &&
		payload.Issuer != "accounts.google.com" {
		return nil, errors.New("invalid issuer")
	}

	email, _ := payload.Claims["email"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	return &model.GooglePayload{
		Email:         email,
		EmailVerified: emailVerified,
		Name:          name,
		Picture:       picture,
		Sub:           payload.Subject,
	}, nil
}
