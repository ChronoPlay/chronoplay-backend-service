package service

import (
	"context"

	"github.com/ChronoPlay/chronoplay-backend-service/dto"
	"github.com/ChronoPlay/chronoplay-backend-service/helpers"
	model "github.com/ChronoPlay/chronoplay-backend-service/model"
	"github.com/ChronoPlay/chronoplay-backend-service/utils"
)

type NotificationService interface {
	SendNotification(ctx context.Context, req dto.SendNotificationRequest) *helpers.CustomError
	MarkNotificationsAsRead(ctx context.Context, req dto.MarkNotificationsAsReadRequest) *helpers.CustomError
	GetNotifications(ctx context.Context, req dto.GetNotificationsRequest) ([]model.Notification, *helpers.CustomError)
	SendDeactivationEmail(ctx context.Context, req dto.SendDeactivationEmailRequest) *helpers.CustomError
}

type notificationService struct {
	notificationRepo model.NotificationRepository
}

func NewNotificationService(notificationRepo model.NotificationRepository) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
	}
}

func (s *notificationService) SendNotification(ctx context.Context, req dto.SendNotificationRequest) *helpers.CustomError {
	var notifications []model.Notification
	for _, userId := range req.UserIds {
		notification := model.Notification{
			UserId:  userId,
			Title:   req.Title,
			Message: req.Message,
			Read:    false,
		}
		notifications = append(notifications, notification)
	}
	err := s.notificationRepo.AddNotifications(ctx, notifications)
	if err != nil {
		return err
	}
	return nil
}

func (s *notificationService) MarkNotificationsAsRead(ctx context.Context, req dto.MarkNotificationsAsReadRequest) *helpers.CustomError {
	if !req.ReadAll && req.NotificationId == 0 {
		return nil
	}
	if req.ReadAll {
		err := s.notificationRepo.MarkNotificationsAsRead(ctx, req.UserId, nil)
		if err != nil {
			return err
		}
		return nil
	}
	err := s.notificationRepo.MarkNotificationsAsRead(ctx, req.UserId, []uint32{req.NotificationId})
	if err != nil {
		return err
	}
	return nil
}

func (s *notificationService) GetNotifications(ctx context.Context, req dto.GetNotificationsRequest) ([]model.Notification, *helpers.CustomError) {
	notifications, err := s.notificationRepo.GetNotificationsByUserId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (s *notificationService) SendDeactivationEmail(ctx context.Context, req dto.SendDeactivationEmailRequest) *helpers.CustomError {
	body := "Dear User,\n\nYou have been Terminated due to insufficient funds for survival tonight. If this is a mistake, please contact support.\n\nBest regards,\nChronoPlay Team"
	err := utils.SendEmail(req.Emails,
		"Survival Failure - Account Deactivation",
		body)
	if err != nil {
		return helpers.System("Failed to send deactivation email: " + err.Error())
	}
	return nil
}
