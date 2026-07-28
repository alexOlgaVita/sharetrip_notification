package domain

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"job4j.ru/sharetrip-notification/internal/dto"
	"job4j.ru/sharetrip-notification/internal/interface"
)

func (u *NotificationUsecase) GetNotification(
	ctx context.Context,
	tx _interface.DBTxer,
	notificationId string,
) (*dto.Notification, error) {

	notification, err := u.NotificationRepo.GetByID(ctx, tx, notificationId)
	if err != nil {
		log.Errorw("s.Repository.Get", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	return notification, nil
}
