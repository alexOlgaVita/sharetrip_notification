package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_interface "job4j.ru/sharetrip-notification/internal/interface"
)

type PoolAdapter struct {
	Pool *pgxpool.Pool
}

func NewPoolAdapter(pool *pgxpool.Pool) *PoolAdapter {
	return &PoolAdapter{Pool: pool}
}

func (a *PoolAdapter) Begin(ctx context.Context) (_interface.DBTxer, error) {
	tx, err := a.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &TxAdapter{tx: tx}, nil
}

func (a *PoolAdapter) Query(ctx context.Context, sql string, args ...any) (_interface.Rows, error) {
	rows, err := a.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &RowsAdapter{rows: rows}, nil
}

func (a *PoolAdapter) QueryRow(ctx context.Context, sql string, args ...any) _interface.Row {
	return &RowAdapter{row: a.Pool.QueryRow(ctx, sql, args...)}
}

func (a *PoolAdapter) Exec(ctx context.Context, sql string, args ...any) (any, error) {
	res, err := a.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *PoolAdapter) Close() error {
	a.Pool.Close()
	return nil
}

func (a *PoolAdapter) Ping(ctx context.Context) error {
	return a.Pool.Ping(ctx)
}

// -------------------------------------------------------------------------
// TxAdapter — реализует _interface.DBTxer для pgx.Tx
// -------------------------------------------------------------------------

type TxAdapter struct {
	tx pgx.Tx
}

func (a *TxAdapter) Query(ctx context.Context, sql string, args ...any) (_interface.Rows, error) {
	rows, err := a.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &RowsAdapter{rows: rows}, nil
}

func (a *TxAdapter) QueryRow(ctx context.Context, sql string, args ...any) _interface.Row {
	return &RowAdapter{row: a.tx.QueryRow(ctx, sql, args...)}
}

func (a *TxAdapter) Exec(ctx context.Context, sql string, args ...any) (any, error) {
	res, err := a.tx.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *TxAdapter) Commit(ctx context.Context) error {
	return a.tx.Commit(ctx)
}

func (a *TxAdapter) Rollback(ctx context.Context) error {
	return a.tx.Rollback(ctx)
}

// -------------------------------------------------------------------------
// RowsAdapter — реализует _interface.Rows для pgx.Rows
// -------------------------------------------------------------------------

type RowsAdapter struct {
	rows pgx.Rows
}

func (a *RowsAdapter) Next() bool {
	return a.rows.Next()
}

func (a *RowsAdapter) Scan(dest ...any) error {
	return a.rows.Scan(dest...)
}

func (a *RowsAdapter) Close() error {
	if a.rows != nil {
		a.rows.Close() // Вызываем метод
	}
	return nil // Явно возвращаем nil, чтобы удовлетворить контракт error
}

func (a *RowsAdapter) Err() error {
	return a.rows.Err()
}

// -------------------------------------------------------------------------
// RowAdapter — реализует _interface.Row для pgx.Row
// -------------------------------------------------------------------------

type RowAdapter struct {
	row pgx.Row
}

func (a *RowAdapter) Scan(dest ...any) error {
	return a.row.Scan(dest...)
}
