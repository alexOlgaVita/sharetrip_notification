package service

import (
	"context"
	"fmt"
	"job4j.ru/sharetrip-notification/internal/interface"

	"go.opentelemetry.io/otel"
	"job4j.ru/sharetrip-notification/internal/domain"
	"job4j.ru/sharetrip-notification/internal/dto"
)

type NotificationService struct {
	Pool _interface.DBQuerier // Используем интерфейс вместо *pgxpool.Pool

	NotificationUsecase *domain.NotificationUsecase
}

func NewNotificationService(
	Pool _interface.DBQuerier,
	NotificationUsecase *domain.NotificationUsecase,
) *NotificationService {
	return &NotificationService{
		Pool:                Pool,
		NotificationUsecase: NotificationUsecase,
	}
}

func (s *NotificationService) CreateNotification(
	ctx context.Context,
	req dto.CreateNotificationRequest,
) (*dto.Notification, error) {
	ctx, span := otel.Tracer("NotificationService").Start(ctx, "NotificationService.CreateNotification")
	defer span.End()

	res, err := tx(ctx, s.Pool, func(tx _interface.DBTxer) (*dto.Notification, error) {

		resp, err := s.NotificationUsecase.CreateNotification(ctx, tx, dto.CreateNotificationRequest{
			RecipientId: req.RecipientId,
			Type:        req.Type,
			Payload:     req.Payload,
		})
		if err != nil {
			return nil, fmt.Errorf("usecase.CreateNotification: %w", err)
		}

		return resp, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return res, nil
}
