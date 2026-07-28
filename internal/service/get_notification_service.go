package service

import (
	"context"
	"fmt"
	"job4j.ru/sharetrip-notification/internal/interface"

	"job4j.ru/sharetrip-notification/internal/dto"
)

func (s *NotificationService) GetNotification(
	ctx context.Context,
	notificationId string,
) (*dto.Notification, error) {

	res, err := tx(ctx, s.Pool, func(tx _interface.DBTxer) (*dto.Notification, error) {
		resp, err := s.NotificationUsecase.GetNotification(ctx, tx, notificationId)
		if err != nil {
			return nil, fmt.Errorf("usecase.GetNotification: %w", err)
		}
		return resp, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return res, nil
}
