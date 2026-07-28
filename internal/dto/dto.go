package dto

type NotificationRequest struct {
	ID          string      `json:"id"`
	RecipientId string      `json:"recipientId"`
	Type        string      `json:"type"`
	Payload     interface{} `json:"payload"`
	Status      string      `json:"status"`
}

type Notification struct {
	ID          string
	RecipientId string
	Type        string
	Payload     interface{}
	Status      string
	CreatedAt   string
}

type CreateNotificationRequest struct {
	RecipientId string      `json:"recipientId"`
	Type        string      `json:"type"`
	Payload     interface{} `json:"payload"`
}

type CreateNotificationResponse struct {
	Notificatio NotificationRequest `json:"Notification"`
}

type UpdateNotificationRequest struct {
	RecipientId string `json:"recipientId"`
	Type        string `json:"type"`
	Payload     string `json:"payload"`

	NotificationID string
	Status         string
}

type GetNotificationResponse struct {
	Notification Notification `json:"Notification"`
}

const (
	NotificationStatusCreated  = "created"
	NotificationStatusPending  = "pending"
	NotificationStatusSent     = "sent"
	NotificationStatusFailed   = "failed"
	NotificationStatusCanceled = "canceled"
)
