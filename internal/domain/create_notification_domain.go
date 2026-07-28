package domain

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"job4j.ru/sharetrip-notification/internal/dto"
	"job4j.ru/sharetrip-notification/internal/interface"
	"job4j.ru/sharetrip-notification/internal/repository"
)

type NotificationUsecase struct {
	// Используем ИНТЕРФЕЙС, а не структуру RepoPg
	NotificationRepo repository.NotificationRepository
}

func (u *NotificationUsecase) CreateNotification(
	ctx context.Context,
	tx _interface.DBTxer,
	req dto.CreateNotificationRequest,
) (*dto.Notification, error) {

	id := uuid.NewString()

	notification, err := u.NotificationRepo.Create(ctx, dto.Notification{
		ID:          id,
		RecipientId: req.RecipientId,
		Type:        req.Type,
		Payload:     req.Payload,
		Status:      dto.NotificationStatusCreated,
	})
	if err != nil {
		return nil, fmt.Errorf("repoNotification.Create: %w", err)
	}

	return notification, nil
}
