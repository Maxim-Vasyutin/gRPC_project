-- +goose Up
CREATE TABLE sessions (
    refresh_token TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ not NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE sessions;
