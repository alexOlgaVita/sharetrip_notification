package api

import (
	"github.com/gofiber/fiber/v2"
	"job4j.ru/sharetrip-notification/internal/dto"
)

type NotificationRequest dto.NotificationRequest

type CreateNotificationRequest dto.CreateNotificationRequest

type CreateNotificationResponse dto.CreateNotificationResponse

func (s *Server) CreateNotification(c *fiber.Ctx) error {

	var req CreateNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	if err := createValidate(&req); err != nil {
		return err
	}

	resp, err := s.NotificationService.CreateNotification(c.Context(),
		dto.CreateNotificationRequest{
			RecipientId: req.RecipientId,
			Type:        req.Type,
			Payload:     req.Payload,
		})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func createValidate(req *CreateNotificationRequest) error {
	if req.RecipientId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "recipientId is required")
	}

	if req.Type == "" {
		return fiber.NewError(fiber.StatusBadRequest, "type is required")
	}

	if req.Payload == "" {
		return fiber.NewError(fiber.StatusBadRequest, "payload is required")
	}

	return nil
}
