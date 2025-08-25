package crons

import (
	"log"
	"time"
)

func SurvivalTaxTask() {
	log.Printf("Survival tax task executed at time: %v", time.Now())
}
