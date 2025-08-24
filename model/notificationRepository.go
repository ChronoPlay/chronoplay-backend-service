package model

import (
	"context"

	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Notification struct {
	Id             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	NotificationId uint32               `bson:"notification_id" json:"notification_id"`
	UserId         uint32               `bson:"user_id" json:"user_id"`
	Title          string               `bson:"title" json:"title"`
	Message        string               `bson:"message" json:"message"`
	Read           bool                 `bson:"read" json:"read"`
	Category       NotificationCategory `bson:"category" json:"category"`
}

type NotificationCategory string

const (
	NotificationCategoryInfo    NotificationCategory = "info"
	NotificationCategoryWarning NotificationCategory = "warning"
	NotificationCategoryAlert   NotificationCategory = "alert"
)

type NotificationRepository interface {
	GetCollection() *mongo.Collection
	AddNotifications(ctx context.Context, notifications []Notification) *helpers.CustomError
	GetNotificationsByUserId(ctx context.Context, userId uint32) ([]Notification, *helpers.CustomError)
	MarkNotificationsAsRead(ctx context.Context, userId uint32, notificationIds []uint32) *helpers.CustomError
}

type mongoNotificationRepo struct {
	collection *mongo.Collection
}

func NewNotificationRepository(col *mongo.Collection) NotificationRepository {
	return &mongoNotificationRepo{collection: col}
}

func (repo *mongoNotificationRepo) GetCollection() *mongo.Collection {
	return repo.collection
}

func (repo *mongoNotificationRepo) AddNotifications(ctx context.Context, notifications []Notification) *helpers.CustomError {
	for i := range notifications {
		nextId, err := GetNextSequence(ctx, repo.collection.Database(), "notificationIds")
		if err != nil {
			return helpers.System("Failed to generate notification ID: " + err.Error())
		}
		notifications[i].NotificationId = uint32(nextId)
	}
	for i := range notifications {
		if notifications[i].Category == "" {
			notifications[i].Category = NotificationCategoryInfo
		}
	}
	docs := make([]interface{}, len(notifications))
	for i := range notifications {
		docs[i] = notifications[i]
	}

	_, err := repo.collection.InsertMany(ctx, docs)
	if err != nil {
		return helpers.System("Failed to add notification: " + err.Error())
	}
	return nil
}

func (repo *mongoNotificationRepo) GetNotificationsByUserId(ctx context.Context, userId uint32) ([]Notification, *helpers.CustomError) {
	var notifications []Notification
	cursor, err := repo.collection.Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		return nil, helpers.System("Failed to fetch notifications: " + err.Error())
	}
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, helpers.System("Failed to decode notifications: " + err.Error())
	}
	return notifications, nil
}

func (repo *mongoNotificationRepo) MarkNotificationsAsRead(ctx context.Context, userId uint32, notificationIds []uint32) *helpers.CustomError {
	var filter bson.M
	if len(notificationIds) == 0 {
		filter = bson.M{"user_id": userId, "read": false}
	} else {
		filter = bson.M{"user_id": userId, "notification_id": bson.M{"$in": notificationIds}, "read": false}
	}
	update := bson.M{
		"$set": bson.M{"read": true},
	}
	_, err := repo.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return helpers.System("Failed to mark notifications as read: " + err.Error())
	}
	return nil
}
