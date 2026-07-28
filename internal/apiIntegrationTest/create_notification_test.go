package apiIntegrationTest_test

import (
	"bytes"
	"encoding/json"
	"io"
	"job4j.ru/sharetrip-notification/internal/dto"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestServer_CreateNotification(t *testing.T) {
	t.Run("success - создание уведомления", func(t *testing.T) {
		payload := dto.CreateNotificationRequest{
			RecipientId: uuid.NewString(),
			Type:        "trip_published",
			Payload:     "{\"trip_id\": \"trip-456\"}",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPost,
			"/notification/",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req, -1)
		require.NoError(t, err)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("close response body: %v", err)
			}
		}()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var got dto.Notification
		err = json.Unmarshal(respBody, &got)
		require.NoError(t, err)
		require.Equal(t,
			dto.CreateNotificationRequest{
				RecipientId: payload.RecipientId,
				Type:        payload.Type,
				Payload:     payload.Payload,
			},
			dto.CreateNotificationRequest{
				RecipientId: got.RecipientId,
				Type:        got.Type,
				Payload:     got.Payload,
			})
	})
}
