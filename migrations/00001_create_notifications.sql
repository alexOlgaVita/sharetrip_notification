-- +goose Up
-- +goose StatementBegin

CREATE TYPE notification_status AS ENUM (
'created',
'pending',
'sent',
'failed',
'canceled'
);

CREATE TABLE IF NOT EXISTS notifications (
    id          UUID PRIMARY KEY,
    user_id     TEXT NOT NULL,
    type TEXT NOT NULL,
    payload     JSONB        NOT NULL DEFAULT '{}',
    status      notification_status  NOT NULL DEFAULT 'created',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
    );

CREATE INDEX idx_notifications_user_id_created_at
    ON notifications (user_id, created_at DESC);

CREATE INDEX idx_notifications_status_created_at
    ON notifications (status, created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS notifications;
DROP TYPE IF EXISTS notification_status;

-- +goose StatementEnd
