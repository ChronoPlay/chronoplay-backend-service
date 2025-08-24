package dto

type MarkNotificationsAsReadRequest struct {
	UserId         uint32 `json:"user_id"`
	NotificationId uint32 `json:"notification_id"`
	ReadAll        bool   `json:"read_all"`
}

type SendNotificationRequest struct {
	UserIds []uint32 `json:"user_id"`
	Title   string   `json:"title"`
	Message string   `json:"message"`
}

type GetNotificationsRequest struct {
	UserId uint32 `json:"user_id"`
}
