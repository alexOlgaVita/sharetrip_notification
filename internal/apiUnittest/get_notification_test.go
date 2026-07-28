package apiUnittest_test

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"job4j.ru/sharetrip-notification/internal/dto"
	"net/http"
	"testing"
)

func TestGetNotificationNotification(t *testing.T) {
	// Инициализируем окружение через наш хелпер
	env := setupTest()

	t.Run("success - получение существующего уведомления", func(t *testing.T) {
		expectedID := uuid.NewString()
		expectedNotification := &dto.Notification{
			ID:          expectedID,
			RecipientId: uuid.NewString(),
			Type:        "email",
			Payload:     "{\"trip_id\": \"trip-456\"}",
			Status:      dto.NotificationStatusCreated,
			CreatedAt:   "2024-01-01T00:00:00Z",
		}

		// Настраиваем мок репозитория
		// Мы используем env.MockTx, так как setupTest уже привязал его к Begin
		env.MockRepo.On("GetByID", mock.Anything, env.MockTx, expectedID).
			Return(expectedNotification, nil).Once()

		url := "/notification/" + expectedID
		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)

		resp, err := env.App.Test(req, -1)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var response dto.GetNotificationResponse
		err = json.Unmarshal(body, &response)
		require.NoError(t, err)

		// Проверка полей (т.к. в DTO структура по значению)
		require.Equal(t, expectedNotification.ID, response.Notification.ID)
	})

	t.Run("error - уведомление не найдено", func(t *testing.T) {
		notFoundID := uuid.NewString()

		// Настраиваем мок на ошибку
		env.MockRepo.On("GetByID", mock.Anything, env.MockTx, notFoundID).
			Return(nil, errors.New("not found")).Once()

		url := "/notification/" + notFoundID
		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)

		resp, err := env.App.Test(req, -1)
		require.NoError(t, err)

		// Проверяем, что API вернул ошибку (4xx или 5xx)
		require.True(t, resp.StatusCode >= 400)
	})
}
