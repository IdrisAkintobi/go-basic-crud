-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    userId UUID NOT NULL,
    deviceId UUID NOT NULL,
    token TEXT UNIQUE NOT NULL,
    userAgent TEXT,
    ipAddress TEXT,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expiresAt TIMESTAMP NOT NULL,
    CONSTRAINT unique_user_device UNIQUE (userId, deviceId)
);

CREATE INDEX idx_sessions_userId ON sessions(userId);
CREATE INDEX idx_sessions_token ON sessions(token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_sessions_userId;
DROP INDEX idx_sessions_token
DROP TABLE sessions;
-- +goose StatementEnd
