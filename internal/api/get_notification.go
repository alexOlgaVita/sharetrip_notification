package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"job4j.ru/sharetrip-notification/internal/domain"
	"job4j.ru/sharetrip-notification/internal/dto"
)

type GetNotificationResponse struct {
	Notification dto.Notification `json:"notification"`
}

func (s *Server) GetNotification(c *fiber.Ctx) error {
	notificationId := c.Params("notificationId")
	if notificationId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "notificationId is required")
	}

	notification, err := s.NotificationService.GetNotification(c.Context(), notificationId)

	if err != nil {
		if errors.Is(err, domain.ErrNotificationNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "notificationId is not found")
		}
		log.Errorw(
			"get notificationId failed",
			"error", err,
			"notification_id", notificationId,
		)
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	res := dto.Notification{
		ID:          notification.ID,
		RecipientId: notification.RecipientId,
		Type:        notification.Type,
		Payload:     notification.Payload,
		Status:      notification.Status,
		CreatedAt:   notification.CreatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(GetNotificationResponse{Notification: res})
}
