package apiUnittest_test

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"job4j.ru/sharetrip-notification/internal/api"
	"job4j.ru/sharetrip-notification/internal/domain"
	"job4j.ru/sharetrip-notification/internal/dto"
	_interface "job4j.ru/sharetrip-notification/internal/interface"
	"job4j.ru/sharetrip-notification/internal/service"
)

type MockRepo struct{ mock.Mock }

func (m *MockRepo) GetByID(ctx context.Context, tx _interface.DBTxer, id string) (*dto.Notification, error) {
	args := m.Called(ctx, tx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Notification), args.Error(1)
}

func (m *MockRepo) Create(ctx context.Context, notification dto.Notification) (*dto.Notification, error) {
	args := m.Called(ctx, notification)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Notification), args.Error(1)
}

func (m *MockRepo) DoPing(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockTx struct{ mock.Mock }

func (m *MockTx) Commit(ctx context.Context) error   { return m.Called(ctx).Error(0) }
func (m *MockTx) Rollback(ctx context.Context) error { return m.Called(ctx).Error(0) }

func (m *MockTx) Exec(ctx context.Context, sql string, args ...any) (any, error) {
	ret := m.Called(ctx, sql, args)
	return ret.Get(0), ret.Error(1)
}

func (m *MockTx) Query(ctx context.Context, sql string, args ...any) (_interface.Rows, error) {
	ret := m.Called(ctx, sql, args)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(_interface.Rows), ret.Error(1)
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...any) _interface.Row {
	ret := m.Called(ctx, sql, args)
	return ret.Get(0).(_interface.Row)
}

type MockPool struct{ mock.Mock }

func (m *MockPool) Begin(ctx context.Context) (_interface.DBTxer, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(_interface.DBTxer), args.Error(1)
}

func (m *MockPool) Query(ctx context.Context, sql string, args ...any) (_interface.Rows, error) {
	ret := m.Called(ctx, sql, args)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(_interface.Rows), ret.Error(1)
}

func (m *MockPool) QueryRow(ctx context.Context, sql string, args ...any) _interface.Row {
	ret := m.Called(ctx, sql, args)
	return ret.Get(0).(_interface.Row)
}

func (m *MockPool) Exec(ctx context.Context, sql string, args ...any) (any, error) {
	ret := m.Called(ctx, sql, args)
	return ret.Get(0), ret.Error(1)
}

func (m *MockPool) Close() error { return m.Called().Error(0) }
func (m *MockPool) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type TestEnv struct {
	App      *fiber.App
	MockRepo *MockRepo
	MockPool *MockPool
	MockTx   *MockTx
}

func setupTest() *TestEnv {
	mockRepo := new(MockRepo)
	mockPool := new(MockPool)
	mockTx := new(MockTx)

	mockPool.On("Begin", mock.Anything).Return(mockTx, nil).Maybe()

	mockTx.On("Commit", mock.Anything).Return(nil).Maybe()
	mockTx.On("Rollback", mock.Anything).Return(nil).Maybe()
	mockTx.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	mockTx.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	mockTx.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	usecase := &domain.NotificationUsecase{
		NotificationRepo: mockRepo,
	}

	svc := service.NewNotificationService(mockPool, usecase)

	server := api.NewServer(mockRepo)
	server.NotificationService = svc

	app := fiber.New()
	server.Route(app.Group(""))

	return &TestEnv{
		App:      app,
		MockRepo: mockRepo,
		MockPool: mockPool,
		MockTx:   mockTx,
	}
}
