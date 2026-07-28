package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"job4j.ru/sharetrip-notification/internal/interface"
)

func tx[T interface{}](
	ctx context.Context,
	pool _interface.DBQuerier,
	block func(tx _interface.DBTxer) (*T, error),
) (*T, error) {

	txBegin, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(txBegin _interface.DBTxer, ctx context.Context) {
		err := txBegin.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			return
		}
	}(txBegin, ctx)

	res, err := block(txBegin)
	if err != nil {
		return nil, fmt.Errorf("repoNotification.transaction block: %w", err)
	}

	if err = txBegin.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return res, nil
}
