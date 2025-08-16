package middleware

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/ChronoPlay/chronoplay-backend-service/utils"
)

type contextKey string

const (
	UserIDKey     contextKey = "userId"
	UserTypeKey   contextKey = "userType"
	TrackingIDKey contextKey = "trackingId"
)

func CustomContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userId, exists := c.Get("UserID")
		if exists {
			ctx = context.WithValue(ctx, UserIDKey, userId)
		}

		userType, exists := c.Get("UserType")
		if exists {
			ctx = context.WithValue(ctx, UserTypeKey, userType)
		}

		trackingId := c.GetHeader("tracking-id")
		if trackingId != "" {
			ctx = context.WithValue(ctx, TrackingIDKey, trackingId)
		}

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func AuthorizeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Authorizing user...")
		token := c.GetHeader("Authorization")
		log.Printf("Authorization token: %s\n", token)
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization token not provided"})
			return
		}

		const bearerPrefix = "Bearer "
		if len(token) > len(bearerPrefix) && token[:len(bearerPrefix)] == bearerPrefix {
			token = token[len(bearerPrefix):]
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token format"})
			return
		}
		log.Printf("Parsed token: %s\n", token)

		userId, err := utils.ParseJwtToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}
		log.Printf("User ID from token: %d\n", userId)

		c.Set("UserID", userId)
		c.Next()
	}
}
