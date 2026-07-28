package _interface

import "context"

type DBTxer interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row

	Exec(ctx context.Context, sql string, args ...any) (any, error)

	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
}

type DBQuerier interface {
	Begin(ctx context.Context) (DBTxer, error)
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, args ...any) (any, error)
	Close() error
	Ping(ctx context.Context) error
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

type Row interface {
	Scan(dest ...any) error // Используйте any (или interface{})
}
