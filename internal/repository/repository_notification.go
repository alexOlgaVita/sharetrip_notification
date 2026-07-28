package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"job4j.ru/sharetrip-notification/internal/dto"
	_interface "job4j.ru/sharetrip-notification/internal/interface"
	"log"
)

//  Repository Interface --- Это КЛЮЧЕВОЙ элемент. Без него подмена невозможна.

type NotificationRepository interface {
	GetByID(ctx context.Context, tx _interface.DBTxer, id string) (*dto.Notification, error)
	Create(ctx context.Context, it dto.Notification) (*dto.Notification, error)
	DoPing(ctx context.Context) error
}

type RepoPg struct {
	pool _interface.DBQuerier
}

func NewRepoPg(pool _interface.DBQuerier) *RepoPg {

	return &RepoPg{pool: pool}
}

func (r *RepoPg) Create(ctx context.Context, it dto.Notification) (*dto.Notification, error) {
	payloadBytes, err := json.Marshal(it.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var result dto.Notification

	err = r.pool.QueryRow(
		ctx,
		`INSERT INTO notifications(id, user_id, type, payload, status) 
         VALUES($1, $2, $3, $4, $5) 
         RETURNING id, user_id, type, payload, status, COALESCE(to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '') AS created_at`, // <--- Важно: возвращаем поля
		it.ID, it.RecipientId, it.Type, payloadBytes, it.Status,
	).Scan(
		&result.ID,
		&result.RecipientId,
		&result.Type,
		&result.Payload, // Здесь pgx сам распарсит jsonb обратно в interface{}
		&result.Status,
		&result.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert and return notification: %w", err)
	}

	return &result, nil
}

func (r *RepoPg) List(ctx context.Context) ([]dto.Notification, error) {
	rows, err := r.pool.Query(ctx, `select id, name from notifications`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing rows: %v", err)
		}
	}()

	var notifications []dto.Notification
	for rows.Next() {
		var item dto.Notification
		if err := rows.Scan(&item.ID, &item.RecipientId, &item.Type, &item.Payload); err != nil {
			return nil, err
		}
		notifications = append(notifications, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *RepoPg) Get(ctx context.Context, notificationId string) (dto.Notification, error) {
	var it dto.Notification
	err := r.pool.QueryRow(
		ctx,
		`select id, user_id, type, payload, status from notifications where id = $1`,
		notificationId,
	).Scan(&it.ID, &it.RecipientId, &it.Type, &it.Payload, &it.Status)

	return it, err
}

func (r *RepoPg) Update(ctx context.Context, name string, newName string) error {
	_, err := r.pool.Exec(
		ctx,
		"UPDATE notifications SET name = $2 WHERE name = $1",
		name, newName,
	)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	return nil
}

func (r *RepoPg) UpdateStatus(ctx context.Context, tx _interface.DBTxer, id string, oldStatus string, newStatus string) error {

	_, err := tx.Exec(ctx, "UPDATE notifications SET status = $2 WHERE id = $1", id, newStatus)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	return nil
}

func (r *RepoPg) Delete(ctx context.Context, name string) error {
	_, err := r.pool.Exec(
		ctx,
		"DELETE FROM notifications WHERE name = $1",
		name,
	)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	return nil
}

func (r *RepoPg) GetCount(ctx context.Context) (string, error) {
	var count string
	err := r.pool.QueryRow(
		ctx,
		`select count(*) from notifications`,
	).Scan(&count)

	return count, err
}

func (r *RepoPg) DoPing(ctx context.Context) error {
	err := r.pool.Ping(ctx)
	return err
}

func (r *RepoPg) GetForUpdateByID(
	ctx context.Context,
	tx _interface.DBTxer,
	id string,
) (*dto.Notification, error) {
	notification := &dto.Notification{}
	err := tx.QueryRow(ctx, "SELECT "+
		"id, "+
		"user_id, "+
		"type, "+
		"payload, "+
		"status, "+
		"COALESCE(to_char(created_at, 'MM-DD-YYYY HH24:MI'), '') AS created_at "+
		"FROM notifications WHERE id = $1 FOR UPDATE", id).Scan(
		&notification.ID,
		&notification.RecipientId,
		&notification.Type,
		&notification.Payload,
		&notification.Status,
		&notification.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotificationNotFound
		}
		return nil, fmt.Errorf("query notification by id %s: %w", id, err)
	}

	return notification, nil
}

func (r *RepoPg) GetByID(
	ctx context.Context,
	tx _interface.DBTxer,
	id string,
) (*dto.Notification, error) {
	notification := &dto.Notification{}

	err := tx.QueryRow(
		ctx,
		`select id, user_id, type, payload, status, 
       COALESCE(to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'), '') from notifications where id = $1 `,
		id).Scan(
		&notification.ID,
		&notification.RecipientId,
		&notification.Type,
		&notification.Payload,
		&notification.Status,
		&notification.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotificationNotFound
		}
		return nil, fmt.Errorf("query notification by id %s: %w", id, err)
	}

	return notification, nil
}
