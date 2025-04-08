package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
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
