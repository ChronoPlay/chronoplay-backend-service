package crons

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/model"
)

func (ctl *cronController) SurvivalTaxTask() {
	if !ctl.cronEnabled {
		log.Println("Cron jobs are disabled. Skipping Survival Tax Task.")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	log.Println("Starting Survival Tax Task...")
	users, err := ctl.userService.GetAllActiveUsers()
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return
	}
	amountForSurvivalTax := float32(50.0)

	deactivatedEmails := []string{}
	for _, user := range users {
		if user.Cash < amountForSurvivalTax {
			user.Cash = 0
			user.Cards = []model.CardOccupied{}
			user.Deactivated = true
			if err != nil {
				log.Printf("Error sending deactivation email to user %d: %v", user.UserId, err)
			}
		} else {
			user.Cash -= amountForSurvivalTax
		}
		err := ctl.userService.UpdateUser(ctx, user)
		if err != nil {
			log.Printf("Error updating user %d: %v", user.UserId, err)
			continue
		}
		if user.Deactivated {
			deactivatedEmails = append(deactivatedEmails, user.Email)
		} else {
			err := ctl.notificationService.SendNotification(ctx, dto.SendNotificationRequest{
				UserIds: []uint32{user.UserId},
				Title:   "Survival Tax Deducted",
				Message: fmt.Sprintf("A survival tax of %.2f has been deducted from your account. Remaining balance: %.2f", amountForSurvivalTax, user.Cash),
			})
			if err != nil {
				log.Printf("Error sending survival tax notification to user %d: %v", user.UserId, err)
			}
		}
		log.Printf("User %d updated successfully. New Cash: %.2f, Deactivated: %v", user.UserId, user.Cash, user.Deactivated)
	}
	if len(deactivatedEmails) > 0 {
		err := ctl.notificationService.SendDeactivationEmail(ctx, dto.SendDeactivationEmailRequest{
			Emails: deactivatedEmails,
		})
		if err != nil {
			log.Printf("Error sending deactivation emails: %v", err)
		} else {
			log.Printf("Deactivation emails sent to: %v", deactivatedEmails)
		}
	}
	log.Println("Survival Tax Task completed.")
}
