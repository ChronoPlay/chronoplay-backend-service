package utils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ParseTime(value primitive.DateTime) time.Time {
	return value.Time()
}
