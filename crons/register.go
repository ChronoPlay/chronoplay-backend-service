package crons

import (
	"log"

	"github.com/robfig/cron/v3"
)

func RunAllCrons() {

	// register survival tax cron
	c := cron.New()

	// Register jobs here
	log.Printf("Registering survival tax cron to run every 3 seconds")
	_, err := c.AddFunc("@every 3s", SurvivalTaxTask)
	if err != nil {
		log.Printf("Error registering survival tax cron: %v", err)
	}
	c.Start()
	log.Println("Cron scheduler started")
}
