package crons

import (
	"log"

	service "github.com/ChronoPlay/chronoplay-backend-service/services"
	"github.com/robfig/cron/v3"
)

type cronController struct {
	userService         service.UserService
	notificationService service.NotificationService
	cronEnabled         bool
}

type CronController interface {
	RunAllCrons()
}

func NewCronController(userService service.UserService, notificationService service.NotificationService, cronEnabled bool) CronController {
	return &cronController{
		userService:         userService,
		notificationService: notificationService,
		cronEnabled:         cronEnabled,
	}
}

func (ctl *cronController) RunAllCrons() {
	// register survival tax cron
	c := cron.New()

	// Register jobs here
	log.Printf("Registering survival tax cron to run every day at midnight")
	_, err := c.AddFunc("0 0 * * *", ctl.SurvivalTaxTask)
	if err != nil {
		log.Printf("Error registering survival tax cron: %v", err)
	}
	c.Start()
	log.Println("Cron scheduler started")
}
